package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTransPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTransPeerCreate,
		Read:   resourceTransPeerRead,
		Update: resourceTransPeerUpdate,
		Delete: resourceTransPeerDelete,

		Schema: map[string]*schema.Schema{
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"nexthop": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"reachable_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTransPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transPeer := &goaviatrix.TransPeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}

	log.Printf("[INFO] Creating Aviatrix transitive peering: %#v", transPeer)

	err := client.CreateTransPeer(transPeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transitive peering: %s", err)
	}
	d.SetId(transPeer.Source + transPeer.Nexthop + transPeer.ReachableCidr)
	//return nil
	return resourceTransPeerRead(d, meta)
}

func resourceTransPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transPeer := &goaviatrix.TransPeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}
	transPeer, err := client.GetTransPeer(transPeer)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transitive peering: %s", err)
	}

	d.Set("source", transPeer.Source)
	d.Set("nexthop", transPeer.Nexthop)
	d.Set("reachable_cidr", transPeer.ReachableCidr)
	return nil
}

func resourceTransPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("aviatrix transitive peering cannot be updated - delete and create new one")
}

func resourceTransPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transPeer := &goaviatrix.TransPeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix transpeer: %#v", transPeer)

	err := client.DeleteTransPeer(transPeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transpeer: %s", err)
	}
	return nil
}
