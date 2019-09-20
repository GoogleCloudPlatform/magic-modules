package google

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform/helper/customdiff"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceBigtableInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableInstanceCreate,
		Read:   resourceBigtableInstanceRead,
		Update: resourceBigtableInstanceUpdate,
		Delete: resourceBigtableInstanceDestroy,

		CustomizeDiff: customdiff.All(
			resourceBigtableInstanceValidateDevelopment,
			resourceBigtableInstanceClusterReorderTypeList,
			resourceResize,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"num_nodes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(3),
						},
						"storage_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SSD",
							ValidateFunc: validation.StringInSlice([]string{"SSD", "HDD"}, false),
						},
					},
				},
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PRODUCTION",
				ValidateFunc: validation.StringInSlice([]string{"DEVELOPMENT", "PRODUCTION"}, false),
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},

			"num_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},

			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},
		},
	}
}

func resourceBigtableInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	conf := &bigtable.InstanceWithClustersConfig{
		InstanceID: d.Get("name").(string),
	}

	displayName, ok := d.GetOk("display_name")
	if !ok {
		displayName = conf.InstanceID
	}
	conf.DisplayName = displayName.(string)

	switch d.Get("instance_type").(string) {
	case "DEVELOPMENT":
		conf.InstanceType = bigtable.DEVELOPMENT
	case "PRODUCTION":
		conf.InstanceType = bigtable.PRODUCTION
	}

	conf.Clusters = expandBigtableClusters(d.Get("cluster").([]interface{}), conf.InstanceID)

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	err = c.CreateInstanceWithClusters(ctx, conf)
	if err != nil {
		return fmt.Errorf("Error creating instance. %s", err)
	}

	d.SetId(conf.InstanceID)

	return resourceBigtableInstanceRead(d, meta)
}

func resourceBigtableInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	instance, err := c.InstanceInfo(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return fmt.Errorf("Error retrieving instance. Could not find %s. %s", d.Id(), err)
	}

	d.Set("project", project)

	clusters := d.Get("cluster").([]interface{})
	clusterState := []map[string]interface{}{}
	for _, cl := range clusters {
		cluster := cl.(map[string]interface{})
		clus, err := c.GetCluster(ctx, instance.Name, cluster["cluster_id"].(string))
		if err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[WARN] Cluster %q not found, not setting it in state", cluster["cluster_id"].(string))
				continue
			}
			return fmt.Errorf("Error retrieving cluster %q: %s", cluster["cluster_id"].(string), err.Error())
		}
		clusterState = append(clusterState, flattenBigtableCluster(clus, cluster["storage_type"].(string)))
	}

	err = d.Set("cluster", clusterState)
	if err != nil {
		return fmt.Errorf("Error setting clusters in state: %s", err.Error())
	}

	d.Set("name", instance.Name)
	d.Set("display_name", instance.DisplayName)

	return nil
}

func resourceBigtableInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}
	defer c.Close()

	clusters, err := c.Clusters(ctx, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving clusters for instance %s", err.Error())
	}

	clusterMap := make(map[string]*bigtable.ClusterInfo, len(clusters))
	for _, cluster := range clusters {
		clusterMap[cluster.Name] = cluster
	}

	clusterTodoMap := make(map[string]interface{}, d.Get("cluster.#").(int))
	var storageTypeShort string
	for _, cluster := range d.Get("cluster").([]interface{}) {
		config := cluster.(map[string]interface{})
		storageTypeShort = config["storage_type"].(string)
		clusterTodoMap[config["cluster_id"].(string)] = config
		cluster_id := config["cluster_id"].(string)
		if cluster, ok := clusterMap[cluster_id]; ok {
			if cluster.ServeNodes != config["num_nodes"].(int) {
				err = c.UpdateCluster(ctx, d.Get("name").(string), cluster.Name, int32(config["num_nodes"].(int)))
				if err != nil {
					return fmt.Errorf("Error updating cluster %s for instance %s", cluster.Name, d.Get("name").(string))
				}
			}
		} else {
			//  new cluster to be added
			var newStorageType bigtable.StorageType
			switch config["storage_type"].(string) {
			case "SSD":
				newStorageType = bigtable.SSD
			case "HDD":
				newStorageType = bigtable.HDD
			}
			conf := &bigtable.ClusterConfig{
				InstanceID:  d.Get("name").(string),
				ClusterID:   cluster_id,
				Zone:        config["zone"].(string),
				NumNodes:    int32(config["num_nodes"].(int)),
				StorageType: newStorageType,
			}
			c.CreateCluster(ctx, conf)
		}
	}

	for name := range clusterMap {
		if _, ok := clusterTodoMap[name]; !ok {
			// cluster is running, but TF will have to remove it
			err = c.DeleteCluster(ctx, d.Get("name").(string), name)
			if err != nil {
				return fmt.Errorf("Error deleting cluster %s for instance %s", name, d.Get("name").(string))
			}
		}
	}
	// build current list of clusters after the change
	clusterState := []map[string]interface{}{}
	clustersOnline, err := c.Clusters(ctx, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving clusters %q: %s", d.Get("name").(string), err.Error())
	}
	for _, ci := range clustersOnline {
		clusterState = append(clusterState, flattenBigtableCluster(ci, storageTypeShort))
	}
	err = d.Set("cluster", clusterState)

	return resourceBigtableInstanceRead(d, meta)
}

func resourceBigtableInstanceDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	name := d.Id()
	err = c.DeleteInstance(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting instance. %s", err)
	}

	d.SetId("")

	return nil
}

func flattenBigtableCluster(c *bigtable.ClusterInfo, storageType string) map[string]interface{} {
	return map[string]interface{}{
		"zone":         c.Zone,
		"num_nodes":    c.ServeNodes,
		"cluster_id":   c.Name,
		"storage_type": storageType,
	}
}

func expandBigtableClusters(clusters []interface{}, instanceID string) []bigtable.ClusterConfig {
	results := make([]bigtable.ClusterConfig, 0, len(clusters))
	for _, c := range clusters {
		cluster := c.(map[string]interface{})
		zone := cluster["zone"].(string)
		var storageType bigtable.StorageType
		switch cluster["storage_type"].(string) {
		case "SSD":
			storageType = bigtable.SSD
		case "HDD":
			storageType = bigtable.HDD
		}
		results = append(results, bigtable.ClusterConfig{
			InstanceID:  instanceID,
			Zone:        zone,
			ClusterID:   cluster["cluster_id"].(string),
			NumNodes:    int32(cluster["num_nodes"].(int)),
			StorageType: storageType,
		})
	}
	return results
}

func resourceBigtableInstanceValidateDevelopment(diff *schema.ResourceDiff, meta interface{}) error {
	if diff.Get("instance_type").(string) != "DEVELOPMENT" {
		return nil
	}
	if diff.Get("cluster.#").(int) != 1 {
		return fmt.Errorf("config is invalid: instance with instance_type=\"DEVELOPMENT\" should have exactly one \"cluster\" block")
	}
	if diff.Get("cluster.0.num_nodes").(int) != 0 {
		return fmt.Errorf("config is invalid: num_nodes cannot be set for instance_type=\"DEVELOPMENT\"")
	}
	return nil
}

func resourceBigtableInstanceClusterReorderTypeList(diff *schema.ResourceDiff, meta interface{}) error {
	old_count, new_count := diff.GetChange("cluster.#")

	// simulate Required:true, MinItems:1, MaxItems:4 for "cluster"
	if new_count.(int) < 1 {
		return fmt.Errorf("config is invalid: Too few cluster blocks: Should have at least 1 \"cluster\" block")
	}
	if new_count.(int) > 4 {
		return fmt.Errorf("config is invalid: Too many cluster blocks: No more than 4 \"cluster\" blocks are allowed")
	}

	var old_ids []string
	clusters := make(map[string]interface{}, new_count.(int))

	for i := 0; i < new_count.(int) || i < old_count.(int); i++ {
		old_id, new_id := diff.GetChange(fmt.Sprintf("cluster.%d.cluster_id", i))
		_, c := diff.GetChange(fmt.Sprintf("cluster.%d", i))
		if old_id != nil && old_id != "" {
			old_ids = append(old_ids, old_id.(string))
		}
		clusters[new_id.(string)] = c
	}

	// reorder clusters according to the old cluster order
	matching_ids := make(map[string]interface{}, new_count.(int))
	var old_cluster_order []interface{}

	for _, id := range old_ids {
		if c, ok := clusters[id]; ok {
			matching_ids[id] = c
			delete(clusters, id)
		}
	}

	for _, id := range old_ids {
		if c, ok := matching_ids[id]; ok {
			old_cluster_order = append(old_cluster_order, c)
		} else {
			// todo empty clusters
			keys := reflect.ValueOf(clusters).MapKeys()
			c := clusters[keys[0].Interface().(string)]
			delete(clusters, keys[0].Interface().(string))
			old_cluster_order = append(old_cluster_order, c)
		}
	}
	for _, value := range clusters {
		old_cluster_order = append(old_cluster_order, value)
	}

	err := diff.SetNew("cluster", old_cluster_order)
	if err != nil {
		return fmt.Errorf("Error setting cluster diff: %s", err)
	}

	return nil
}

func resourceResize(diff *schema.ResourceDiff, meta interface{}) error {
	old_count, new_count := diff.GetChange("cluster.#")
	if old_count.(int) != new_count.(int) {
		log.Printf("[DEBUG] resize from %d to %d", old_count, new_count)
	}

	for i := 0; i < old_count.(int) || i < new_count.(int); i++ {
		cluster_old, cluster_new := diff.GetChange(fmt.Sprintf("cluster.%d", i))
		cluster_old_map := cluster_old.(map[string]interface{})
		cluster_new_map := cluster_new.(map[string]interface{})
		if cluster_old_map["cluster_id"] == nil {
			log.Printf("[DEBUG] new [%s] will be added", cluster_new_map["cluster_id"].(string))
		} else if cluster_new_map["cluster_id"] == "" {
			log.Printf("[DEBUG] old [%s] will be removed", cluster_old_map["cluster_id"].(string))
		} else {
			log.Printf("[DEBUG] changes from %s to %s", cluster_old_map["cluster_id"].(string), cluster_new_map["cluster_id"].(string))
			if cluster_old_map["storage_type"].(string) != cluster_new_map["storage_type"].(string) {
				diff.ForceNew(fmt.Sprintf("cluster.%d.storage_type", i))
			}
			if cluster_old_map["zone"].(string) != cluster_new_map["zone"].(string) {
				diff.ForceNew(fmt.Sprintf("cluster.%d.zone", i))
			}
			if cluster_old_map["cluster_id"].(string) != cluster_new_map["cluster_id"].(string) {
				diff.ForceNew(fmt.Sprintf("cluster.%d.cluster_id", i))
			}
		}

	}
	return nil
}
