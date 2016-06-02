package icf

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joeswaminathan/icf-sdk-go/icf"
	"cto-github.cisco.com/jswamina/kvs_infra/src/infra/log"
	"strings"
	"time"
)

func resourceIcfInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceIcfInstanceCreate,
		Read:   resourceIcfInstanceRead,
		Update: resourceIcfInstanceUpdate,
		Delete: resourceIcfInstanceDelete,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"catalog": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"provider_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"vdc": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"ICF_VDC",
				}, nil),
			},

			"network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enterprise_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceIcfInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	c := meta.(*icf.Client)

	instance := &icf.Instance{
		Vdc:            d.Get("vdc").(string),
		Catalog:        d.Get("catalog").(string),
		ProviderAccess: d.Get("provider_access").(bool),
		Nics: []icf.InstanceNicInfo{
			{
				Index:   1,
				Dhcp:    false,
				Network: d.Get("network").(string),
			},
		},
	}

	instance, err = c.CreateInstance(instance)
	if err != nil {
		log.Printf("[ERROR] Creating Instance %v", err)
		return
	}

	log.Printf("[INFO] Instance ID: %s", instance.Oid)

	// Wait for the instance to become running so we can get some attributes
	// that aren't available until later.
	log.Printf("[DEBUG] Waiting for instance (%s) to become running",
		instance.Oid)


	stateConf := &resource.StateChangeConf{
		Pending:    []string{icf.StatusCreateInProgress},
		Target:     []string{icf.StatusSuccess, icf.StatusCreateFailed},
		Refresh:    InstanceStateRefreshFunc(c, instance.Oid),
		Timeout:    20 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	instanceRaw, err := stateConf.WaitForState()
	if err != nil {
		log.Printf("[ERROR] Wating for status %v", err)
		err = fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			instance.Oid, err)
		return
	}

	instance = instanceRaw.(*icf.Instance)

	if instance.Status == icf.StatusCreateFailed {
		err = fmt.Errorf("Error creating instance")
		return
	}
	// Store the resulting ID so we can look this up later
	d.SetId(instance.Oid)

	return resourceIcfInstanceRead(d, meta)
}

func resourceIcfInstanceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*icf.Client)

	instance, err := c.GetInstance(d.Id())
	if err != nil {
		errs := fmt.Sprintf("%v", err)
		if strings.Contains(errs, "404") || strings.Contains(errs, "400") {
			return nil
		}
		return fmt.Errorf("Error reading ICF instannce %s error %v", d.Id(), err)
	}
	//d.Set("public_ip", instance.PublicIp)
	d.Set("public_ip", instance.Nics[0].Ip)
	//d.Set("private_ip", instance.PrivateIp)
	d.Set("private_ip", instance.Nics[0].Ip)
	d.Set("enterprise_ip", instance.Nics[0].Ip)
	d.Set("name", instance.Name)

	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": instance.Nics[0].Ip,
	})

	return nil
}

func resourceIcfInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIcfInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*icf.Client)

	if err := c.DeleteInstance(d.Id()); err != nil {
		errs := fmt.Sprintf("%v", err)
		if strings.Contains(errs, "404") || strings.Contains(errs, "400") {
			return nil
		}
		return err
	}

	log.Printf("[DEBUG] Waiting for instance (%s) to become removed",
		d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{icf.StatusDeleteInProgress},
		Target:     []string{icf.StatusDeleted},
		Refresh:    InstanceStateRefreshFunc(c, d.Id()),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		log.Printf("[ERROR] Wating for status %v", err)
		err = fmt.Errorf(
			"Error waiting for instance (%s) to become deleted: %s",
			d.Id(), err)
		return err
	}
	d.SetId("")
	return nil
}

func InstanceStateRefreshFunc(c *icf.Client, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := c.GetInstance(instanceID)
		if err != nil {
			status := ""
			errs := fmt.Sprintf("%v", err)
			if strings.Contains(errs, "404") || strings.Contains(errs, "400") {
				status = icf.StatusDeleted
				err = nil
				return instance, status, err
			}
			log.Printf("[ERROR] InstanceStateRefreshFunc : Error = ", err)
			return instance, status, err
		}

		log.Printf("[INFO] InstanceStateRefreshFunc : Status = ", instance.Status)
		return instance, instance.Status, nil
	}
}
