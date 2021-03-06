// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.
package ovirt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ovirtsdk4 "gopkg.in/imjoey/go-ovirt.v4"
)

func TestAccOvirtVM_basic(t *testing.T) {
	var vm ovirtsdk4.Vm
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: "ovirt_vm.vm",
		CheckDestroy:  testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVMBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVMExists("ovirt_vm.vm", &vm),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "name", "testAccVMBasic"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "status", "up"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.#", "1"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.nic_configuration.#", "2"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.host_name", "vm-basic-1"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.user_name", "root"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.custom_script", "echo hello"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.dns_search", "university.edu"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.dns_servers", "8.8.8.8 8.8.4.4"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.nic_configuration.0.label", "eth0"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.nic_configuration.0.address", "10.1.60.60"),
				),
			},
			{
				Config: testAccVMBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVMExists("ovirt_vm.vm", &vm),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "name", "testAccVMBasic"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "status", "up"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.#", "1"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.nic_configuration.#", "1"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.host_name", "vm-basic-updated"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.user_name", "root"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.custom_script", "echo hello2"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.dns_search", "university.edu"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.dns_servers", "8.8.8.8"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.authorized_ssh_key", ""),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.nic_configuration.0.label", "eth0"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "initialization.0.nic_configuration.0.address", "10.1.60.66"),
				),
			},
		},
	})
}

func TestAccOvirtVM_attachedDisk(t *testing.T) {
	var vm ovirtsdk4.Vm
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: "ovirt_vm.vm",
		CheckDestroy:  testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVMAttachedDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVMExists("ovirt_vm.vm", &vm),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "name", "testAccVMAttachedDisk"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "status", "up"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "attached_disk.#", "1"),
				),
			},
			{
				Config: testAccVMAttachedDiskUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVMExists("ovirt_vm.vm", &vm),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "name", "testAccVMAttachedDisk"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "status", "up"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "attached_disk.#", "2"),
				),
			},
		},
	})
}

func TestAccOvirtVM_vnic(t *testing.T) {
	var vm ovirtsdk4.Vm
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: "ovirt_vm.vm",
		CheckDestroy:  testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVMVnic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVMExists("ovirt_vm.vm", &vm),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "name", "testAccVMVnic"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "status", "up"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "vnic.#", "2"),
				),
			},
			{
				Config: testAccVMVnicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVMExists("ovirt_vm.vm", &vm),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "name", "testAccVMVnic"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "status", "up"),
					resource.TestCheckResourceAttr("ovirt_vm.vm", "vnic.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVMDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovirt_vm" {
			continue
		}
		getResp, err := conn.SystemService().VmsService().
			VmService(rs.Primary.ID).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				continue
			}
			return err
		}
		if _, ok := getResp.Vm(); ok {
			return fmt.Errorf("VM %s still exist", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOvirtVMExists(n string, v *ovirtsdk4.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No VM ID is set")
		}
		conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
		getResp, err := conn.SystemService().VmsService().
			VmService(rs.Primary.ID).
			Get().
			Send()
		if err != nil {
			return err
		}
		vm, ok := getResp.Vm()
		if ok {
			*v = *vm
			return nil
		}
		return fmt.Errorf("VM %s not exist", rs.Primary.ID)
	}
}

const testAccVMBasic = `
resource "ovirt_vm" "vm" {
	name        = "testAccVMBasic"
	cluster_id  = "${data.ovirt_clusters.default.clusters.0.id}"

	initialization {
		host_name = "vm-basic-1"
		timezone = "Asia/Shanghai"
		user_name = "root"
		custom_script = "echo hello"
		dns_search = "university.edu"
		dns_servers = "8.8.8.8 8.8.4.4"
		authorized_ssh_key = "${file(pathexpand("~/.ssh/id_rsa.pub"))}"
		nic_configuration {
			label       = "eth0"
			boot_proto  = "static"
			address  	= "10.1.60.60"
			gateway     = "10.1.60.1"
			netmask 	= "255.255.255.0"
		}
		nic_configuration {
			label       = "eth1"
			boot_proto  = "static"
			address  	= "10.1.60.61"
			gateway     = "10.1.60.1"
			netmask 	= "255.255.255.0"
		}
	}

	attached_disk {
		disk_id = "${ovirt_disk.vm_disk.id}"
		bootable = true
		interface = "virtio"
	}
}

resource "ovirt_disk" "vm_disk" {
	name              = "vm_disk"
	alias             = "vm_disk"
	size              = 23687091200
	format            = "cow"
	storage_domain_id = "${data.ovirt_storagedomains.data.storagedomains.0.id}"
	sparse            = true
}

data "ovirt_storagedomains" "data" {
	name_regex = "^data"
	search = {
	  criteria       = "name = data and datacenter = ${data.ovirt_datacenters.default.datacenters.0.name}"
	  case_sensitive = false
	}
}

data "ovirt_datacenters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

data "ovirt_clusters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

`

const testAccVMBasicUpdate = `
resource "ovirt_vm" "vm" {
	name        = "testAccVMBasic"
	cluster_id  = "${data.ovirt_clusters.default.clusters.0.id}"

	initialization {
		host_name = "vm-basic-updated"
		timezone = "Asia/Shanghai"
		user_name = "root"
		custom_script = "echo hello2"
		dns_search = "university.edu"
		dns_servers = "8.8.8.8"
		nic_configuration {
			label       = "eth0"
			boot_proto  = "static"
			address  	= "10.1.60.66"
			gateway     = "10.1.60.1"
			netmask 	= "255.255.255.0"
		}
	}

	attached_disk {
		disk_id = "${ovirt_disk.vm_disk.id}"
		bootable = true
		interface = "virtio"
	}

}

resource "ovirt_disk" "vm_disk" {
	name              = "vm_disk"
	alias             = "vm_disk"
	size              = 23687091200
	format            = "cow"
	storage_domain_id = "${data.ovirt_storagedomains.data.storagedomains.0.id}"
	sparse            = true
}

data "ovirt_storagedomains" "data" {
	name_regex = "^data"
	search = {
	  criteria       = "name = data and datacenter = ${data.ovirt_datacenters.default.datacenters.0.name}"
	  case_sensitive = false
	}
}

data "ovirt_datacenters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

data "ovirt_clusters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

`

const testAccVMAttachedDisk = `
resource "ovirt_vm" "vm" {
	name        = "testAccVMAttachedDisk"
	cluster_id  = "${data.ovirt_clusters.default.clusters.0.id}"

	attached_disk {
		disk_id = "${ovirt_disk.vm_disk.id}"
		bootable = true
		interface = "virtio"
	}
}

resource "ovirt_disk" "vm_disk" {
	name              = "vm_disk"
	alias             = "vm_disk"
	size              = 23687091200
	format            = "cow"
	storage_domain_id = "${data.ovirt_storagedomains.data.storagedomains.0.id}"
	sparse            = true
}

data "ovirt_storagedomains" "data" {
	name_regex = "^data"
	search = {
	  criteria       = "name = data and datacenter = ${data.ovirt_datacenters.default.datacenters.0.name}"
	  case_sensitive = false
	}
}

data "ovirt_datacenters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

data "ovirt_clusters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

`

const testAccVMAttachedDiskUpdate = `
resource "ovirt_vm" "vm" {
	name        = "testAccVMAttachedDisk"
	cluster_id  = "${data.ovirt_clusters.default.clusters.0.id}"

	attached_disk {
		disk_id = "${ovirt_disk.vm_disk.id}"
		bootable = true
		interface = "virtio"
	}

	attached_disk {
		disk_id = "${ovirt_disk.vm_disk_2.id}"
		bootable = false
		interface = "virtio"
	}
}

resource "ovirt_disk" "vm_disk" {
	name              = "vm_disk"
	alias             = "vm_disk"
	size              = 23687091200
	format            = "cow"
	storage_domain_id = "${data.ovirt_storagedomains.data.storagedomains.0.id}"
	sparse            = true
}

resource "ovirt_disk" "vm_disk_2" {
	name              = "vm_disk_2"
	alias             = "vm_disk_2"
	size              = 23687091200
	format            = "cow"
	storage_domain_id = "${data.ovirt_storagedomains.data.storagedomains.0.id}"
	sparse            = true
}

data "ovirt_storagedomains" "data" {
	name_regex = "^data"
	search = {
	  criteria       = "name = data and datacenter = ${data.ovirt_datacenters.default.datacenters.0.name}"
	  case_sensitive = false
	}
}

data "ovirt_datacenters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

data "ovirt_clusters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}
`

const testAccVMVnic = `
resource "ovirt_vm" "vm" {
	name        = "testAccVMVnic"
	cluster_id  = "${data.ovirt_clusters.default.clusters.0.id}"

	vnic {
		name  			= "nic1"
		vnic_profile_id = "${ovirt_vnic_profile.vm_vnic_profile.id}"
	}

	vnic {
		name  			= "nic2"
		vnic_profile_id = "${ovirt_vnic_profile.vm_vnic_profile.id}"
	}

	attached_disk {
		disk_id = "${ovirt_disk.vm_disk.id}"
		bootable = true
		interface = "virtio"
	}
}

resource "ovirt_disk" "vm_disk" {
	name              = "vm_disk"
	alias             = "vm_disk"
	size              = 23687091200
	format            = "cow"
	storage_domain_id = "${data.ovirt_storagedomains.data.storagedomains.0.id}"
	sparse            = true
}

data "ovirt_storagedomains" "data" {
	name_regex = "^data"
	search = {
	  criteria       = "name = data and datacenter = ${data.ovirt_datacenters.default.datacenters.0.name}"
	  case_sensitive = false
	}
}

data "ovirt_datacenters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

data "ovirt_clusters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

resource "ovirt_vnic_profile" "vm_vnic_profile" {
	name        	= "vm_vnic_profile"
	network_id  	= "${data.ovirt_networks.ovirtmgmt.networks.0.id}"
	migratable  	= false
	port_mirroring 	= false
}

data "ovirt_networks" "ovirtmgmt" {
	search = {
	  criteria       = "datacenter = Default and name = ovirtmgmt"
	  max            = 1
	  case_sensitive = false
	}
}

`

const testAccVMVnicUpdate = `
resource "ovirt_vm" "vm" {
	name        = "testAccVMVnic"
	cluster_id  = "${data.ovirt_clusters.default.clusters.0.id}"

	vnic {
		name  			= "nic11"
		vnic_profile_id = "${ovirt_vnic_profile.vm_vnic_profile.id}"
	}

	attached_disk {
		disk_id = "${ovirt_disk.vm_disk.id}"
		bootable = true
		interface = "virtio"
	}
}

resource "ovirt_disk" "vm_disk" {
	name              = "vm_disk"
	alias             = "vm_disk"
	size              = 23687091200
	format            = "cow"
	storage_domain_id = "${data.ovirt_storagedomains.data.storagedomains.0.id}"
	sparse            = true
}

data "ovirt_storagedomains" "data" {
	name_regex = "^data"
	search = {
	  criteria       = "name = data and datacenter = ${data.ovirt_datacenters.default.datacenters.0.name}"
	  case_sensitive = false
	}
}

data "ovirt_datacenters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

data "ovirt_clusters" "default" {
	search = {
		criteria       = "name = Default"
		max            = 1
		case_sensitive = false
	}
}

resource "ovirt_vnic_profile" "vm_vnic_profile" {
	name        	= "vm_vnic_profile"
	network_id  	= "${data.ovirt_networks.ovirtmgmt.networks.0.id}"
	migratable  	= false
	port_mirroring 	= false
}

data "ovirt_networks" "ovirtmgmt" {
	search = {
	  criteria       = "datacenter = Default and name = ovirtmgmt"
	  max            = 1
	  case_sensitive = false
	}
}

`
