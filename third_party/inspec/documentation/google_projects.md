### Test that there are no more than a specified number of projects available for the project

    describe google_projects do
      its('count') { should be <= 100}
    end

### Test that an expected named project is available

    describe google_projects do
      its('project_names'){ should include "GCP Project Name" }
    end

### Test that an expected project identifier is available

    describe google_projects do
      its('project_ids'){ should include "gcp_project_id" }
    end
    
### Test that an expected project number is available

    describe google_projects do
      its('project_numbers'){ should include 1122334455 }
    end    

### Test that a particular subset of projects with id 'prod*' are in ACTIVE lifecycle state

    google_projects.where(project_id: /^prod/).project_ids.each do |gcp_project_id|
      describe google_project(project: gcp_project_id) do
        it { should exist }
        its('lifecycle_state') { should eq "ACTIVE" }
      end
    end

### Test that a particular subset of ACTIVE projects with id 'prod*' exist

    google_projects.where(project_id: /^prod/, lifecycle_state: 'ACTIVE').project_ids.each do |gcp_project_id|
      describe google_project(project: gcp_project_id) do
        it { should exist }
      end
    end