# Ansible PR Script
This is a Bash script that creates a series of PRs from  the `build/ansible`
folder to ansible/ansible upstream repo.

Pull Requests are made from the origin remote to Ansible core.
By default, the origin remote points to MM's version of Ansible

## Requirements
* `hub` CLI
* The origin fork on 'build/ansible' points to your fork of Ansible.
  All PRs to upstream Ansible will come from the owner of this origin
  remote fork.
