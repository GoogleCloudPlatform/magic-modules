package google

import (
	"fmt"
	"log"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/googleapi"
)

func resourceComputeRouterNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRouterNatCreate,
		Read:   resourceComputeRouterNatRead,
		Delete: resourceComputeRouterNatDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeRouterNatImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateRFC1035Name(6, 30),
			},

			"router": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_ip_allocate_option": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_ips": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_subnetwork_ip_ranges_to_nat": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnetwork": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"source_ip_ranges_to_nat": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"secondary_ip_range_names": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"min_ports_per_vm": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"udp_idle_timeout_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"icmp_idle_timeout_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tcp_established_idle_timeout_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tcp_transitory_idle_timeout_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeRouterNatCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	routerName := d.Get("router").(string)
	natName := d.Get("name").(string)

	routerLock := getRouterLockName(region, routerName)
	mutexKV.Lock(routerLock)
	defer mutexKV.Unlock(routerLock)

	routersService := config.clientComputeBeta.Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router nat %s because its router %s/%s is gone", natName, region, routerName)
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error Reading router %s/%s: %s", region, routerName, err)
	}

	nats := router.Nats
	for _, nat := range nats {
		if nat.Name == natName {
			d.SetId("")
			return fmt.Errorf("Router %s has nat %s already", routerName, natName)
		}
	}

	nat := &computeBeta.RouterNat{Name: natName}

	if v, ok := d.GetOk("nat_ip_allocate_option"); ok {
		nat.NatIpAllocateOption = v.(string)
	}

	if v, ok := d.GetOk("source_subnetwork_ip_ranges_to_nat"); ok {
		nat.SourceSubnetworkIpRangesToNat = v.(string)
	}

	if v, ok := d.GetOk("min_ports_per_vm"); ok {
		nat.MinPortsPerVm = int64(v.(int))
	}

	if v, ok := d.GetOk("udp_idle_timeout_sec"); ok {
		nat.UdpIdleTimeoutSec = int64(v.(int))
	}

	if v, ok := d.GetOk("icmp_idle_timeout_sec"); ok {
		nat.IcmpIdleTimeoutSec = int64(v.(int))
	}

	if v, ok := d.GetOk("tcp_established_idle_timeout_sec"); ok {
		nat.TcpEstablishedIdleTimeoutSec = int64(v.(int))
	}

	if v, ok := d.GetOk("tcp_transitory_idle_timeout_sec"); ok {
		nat.TcpTransitoryIdleTimeoutSec = int64(v.(int))
	}

	log.Printf("[INFO] Adding nat %s", natName)
	nats = append(nats, nat)
	patchRouter := &computeBeta.Router{
		Nats: nats,
	}

	log.Printf("[DEBUG] Updating router %s/%s with nats: %+v", region, routerName, nats)
	op, err := routersService.Patch(project, region, router.Name, patchRouter).Do()
	if err != nil {
		return fmt.Errorf("Error patching router %s/%s: %s", region, routerName, err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", region, routerName, natName))
	err = computeBetaOperationWaitTime(config.clientCompute, op, project, "Patching router", 4)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error waiting to patch router %s/%s: %s", region, routerName, err)
	}

	return resourceComputeRouterNatRead(d, meta)
}

func resourceComputeRouterNatRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	routerName := d.Get("router").(string)
	natName := d.Get("name").(string)

	routersService := config.clientComputeBeta.Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router nat %s because its router %s/%s is gone", natName, region, routerName)
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error Reading router %s/%s: %s", region, routerName, err)
	}

	for _, nat := range router.Nats {

		if nat.Name == natName {
			d.SetId(fmt.Sprintf("%s/%s/%s", region, routerName, natName))
			d.Set("nat_ip_allocate_option", nat.NatIpAllocateOption)
			d.Set("source_subnetwork_ip_ranges_to_nat", nat.SourceSubnetworkIpRangesToNat)
			d.Set("min_ports_per_vm", nat.MinPortsPerVm)
			d.Set("udp_idle_timeout_sec", nat.UdpIdleTimeoutSec)
			d.Set("icmp_idle_timeout_sec", nat.IcmpIdleTimeoutSec)
			d.Set("tcp_established_idle_timeout_sec", nat.TcpEstablishedIdleTimeoutSec)
			d.Set("tcp_transitory_idle_timeout_sec", nat.TcpTransitoryIdleTimeoutSec)
			d.Set("region", region)
			d.Set("project", project)
			return nil
		}
	}

	log.Printf("[WARN] Removing router nat %s/%s/%s because it is gone", region, routerName, natName)
	d.SetId("")
	return nil
}

func resourceComputeRouterNatDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	routerName := d.Get("router").(string)
	natName := d.Get("name").(string)

	routerLock := getRouterLockName(region, routerName)
	mutexKV.Lock(routerLock)
	defer mutexKV.Unlock(routerLock)

	routersService := config.clientComputeBeta.Routers
	router, err := routersService.Get(project, region, routerName).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing router nat %s because its router %s/%s is gone", natName, region, routerName)

			return nil
		}

		return fmt.Errorf("Error Reading Router %s: %s", routerName, err)
	}

	var newNats []*computeBeta.RouterNat = make([]*computeBeta.RouterNat, 0, len(router.BgpPeers))
	for _, nat := range router.Nats {
		if nat.Name == natName {
			continue
		} else {
			newNats = append(newNats, nat)
		}
	}

	if len(newNats) == len(router.Nats) {
		log.Printf("[DEBUG] Router %s/%s had no nat %s already", region, routerName, natName)
		d.SetId("")
		return nil
	}

	log.Printf(
		"[INFO] Removing nat %s from router %s/%s", natName, region, routerName)
	patchRouter := &computeBeta.Router{
		Nats: newNats,
	}

	if len(newNats) == 0 {
		patchRouter.ForceSendFields = append(patchRouter.ForceSendFields, "Nats")
	}

	log.Printf("[DEBUG] Updating router %s/%s with nats: %+v", region, routerName, newNats)
	op, err := routersService.Patch(project, region, router.Name, patchRouter).Do()
	if err != nil {
		return fmt.Errorf("Error patching router %s/%s: %s", region, routerName, err)
	}

	err = computeBetaOperationWaitTime(config.clientCompute, op, project, "Patching router", 4)
	if err != nil {
		return fmt.Errorf("Error waiting to patch router %s/%s: %s", region, routerName, err)
	}

	d.SetId("")
	return nil
}

func resourceComputeRouterNatImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid router nat specifier. Expecting {region}/{router}/{nat}")
	}

	d.Set("region", parts[0])
	d.Set("router", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}
