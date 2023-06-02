package esxi

import (
	"github.com/hashicorp/terraform"
)

func resourceGuest() *schema.Resource {
	return &schema.Resource{
		Create: resourceGuestCreate,
		// Read:   resourceGuestRead,
		// Update: resourceGuestUpdate,
		// Delete: resourceGuestDelete,
		Schema: map[string]*schema.Schema{
			"esxi_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"guest_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"memsize": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1024,
			},
			"numvcpus": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"guestos": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "otherGuest",
			},
			"network_interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_network": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"mac_address": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceGuestCreate(d *schema.ResourceData, m interface{}) error {
	// Aquí tendrás que hacer la llamada al API de vSphere o ESXi para crear la VM
	// Aquí es donde se deberían leer los valores de tu esquema y usarlos para crear la VM

	esxiHostname := d.Get("esxi_hostname").(string)
	esxiUsername := d.Get("esxi_username").(string)
	esxiPassword := d.Get("esxi_password").(string)
	guestName := d.Get("guest_name").(string)
	memSize := int64(d.Get("memsize").(int))
	numVCPUs := int32(d.Get("numvcpus").(int))
	guestOS := d.Get("guestos").(string)

	err := createVM(esxiHostname, esxiUsername, esxiPassword, guestName, memSize, numVCPUs, guestOS)
	if err != nil {
		return err
	}

	// Aquí es donde necesitarías obtener el ID de la VM creada para usarla como ID del recurso Terraform.
	// Esto se supone que se hace a través de la API de ESXi.
	// Como es solo un ejemplo, establecemos el ID del recurso como el nombre del invitado.
	d.SetId(guestName)

	return nil //resourceGuestRead(d, m)
}
