package google

import (
	"fmt"
	"log"

	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	compute "google.golang.org/api/compute/v0.beta"
)

func resourceComputeRegionSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionSecurityPoliciesCreate,
		Read:   resourceComputeRegionSecurityPoliciesRead,
		Update: resourceComputeRegionSecurityPoliciesUpdate,
		Delete: resourceComputeRegionSecurityPoliciesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeRegionSecurityPoliciesImporter,
		},
		CustomizeDiff: rulesCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCEName,
				Description:  `The name of the security policy.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `An optional description of this security policy. Max size is 2048.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type indicates the intended use of the security policy. CLOUD_ARMOR - Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers. CLOUD_ARMOR_EDGE - Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.`,
				//Change
				//SetFeature: Cloud Armor for NLB/VMs APIs (item c) -> insert CLOUD_ARMOR_NETWORK on enum
				ValidateFunc: validation.StringInSlice([]string{"CLOUD_ARMOR", "CLOUD_ARMOR_EDGE", "CLOUD_ARMOR_INTERNAL_SERVICE", "CLOUD_ARMOR_NETWORK"}, false),
			},

			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // If no rules are set, a default rule is added
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Action to take when match matches the request.`,
						},

						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `An unique positive integer indicating the priority of evaluation for a rule. Rules are evaluated from highest priority (lowest numerically) to lowest priority (highest numerically) in order.`,
						},

						"match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"src_ip_ranges": {
													Type:        schema.TypeSet,
													Required:    true,
													MinItems:    1,
													MaxItems:    10,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: `Set of IP addresses or ranges (IPV4 or IPV6) in CIDR notation to match against inbound traffic. There is a limit of 10 IP ranges per rule. A value of '*' matches all IPs (can be used to override the default behavior).`,
												},
											},
										},
										Description: `The configuration options available when specifying versioned_expr. This field must be specified if versioned_expr is specified and cannot be specified if versioned_expr is not specified.`,
									},

									"versioned_expr": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validation.StringInSlice([]string{"SRC_IPS_V1"}, false),
										Description:  `Predefined rule expression. If this field is specified, config must also be specified. Available options:   SRC_IPS_V1: Must specify the corresponding src_ip_ranges field in config.`,
									},

									"expr": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"expression": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Textual representation of an expression in Common Expression Language syntax. The application context of the containing message determines which well-known feature set of CEL is supported.`,
												},
												// These fields are not yet supported (Issue hashicorp/terraform-provider-google#4497: mbang)
												// "title": {
												// 	Type:     schema.TypeString,
												// 	Optional: true,
												// },
												// "description": {
												// 	Type:     schema.TypeString,
												// 	Optional: true,
												// },
												//"location": {
												//	Type:     schema.TypeString,
												//	Optional: true,
												//},
											},
										},
										Description: `User defined CEVAL expression. A CEVAL expression is used to specify match criteria such as origin.ip, source.region_code and contents in the request header.`,
									},
								},
							},
							Description: `A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding action is enforced.`,
						},

						//MarkWaf
						"preconfigured_waf_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclusion": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"request_header": resourceComputeSecurityPolicyRulePreconfiguredWafConfigExclusionFieldParamsSchema(
													`Request header whose value will be excluded from inspection during preconfigured WAF evaluation.`,
												),

												"request_cookie": resourceComputeSecurityPolicyRulePreconfiguredWafConfigExclusionFieldParamsSchema(
													`Request cookie whose value will be excluded from inspection during preconfigured WAF evaluation.`,
												),

												"request_uri": resourceComputeSecurityPolicyRulePreconfiguredWafConfigExclusionFieldParamsSchema(
													`Request URI from the request line to be excluded from inspection during preconfigured WAF evaluation. When specifying this field, the query or fragment part should be excluded.`,
												),

												"request_query_param": resourceComputeSecurityPolicyRulePreconfiguredWafConfigExclusionFieldParamsSchema(
													`Request query parameter whose value will be excluded from inspection during preconfigured WAF evaluation.  Note that the parameter can be in the query string or in the POST body.`,
												),

												"target_rule_set": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Target WAF rule set to apply the preconfigured WAF exclusion.`,
												},

												"target_rule_ids": {
													Type:        schema.TypeSet,
													Optional:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: `A list of target rule IDs under the WAF rule set to apply the preconfigured WAF exclusion. If omitted, it refers to all the rule IDs under the WAF rule set.`,
												},
											},
										},
										Description: `An exclusion to apply during preconfigured WAF evaluation.`,
									},
								},
							},
							Description: `Preconfigured WAF configuration to be applied for the rule. If the rule does not evaluate preconfigured WAF rules, i.e., if evaluatePreconfiguredWaf() is not used, this field will have no effect.`,
						},

						"description": {
							Type:        schema.TypeString,
							Default:     "",
							Optional:    true,
							Description: `An optional description of this rule. Max size is 64.`,
						},

						"preview": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `When set to true, the action specified above is not enforced. Stackdriver logs for requests that trigger a preview action are annotated as such.`,
						},

						"rate_limit_options": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Rate limit threshold for this security policy. Must be specified if the action is "rate_based_ban" or "throttle". Cannot be specified for any other actions.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rate_limit_threshold": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `Threshold at which to begin ratelimiting.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Number of HTTP(S) requests for calculating the threshold.`,
												},

												"interval_sec": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Interval over which the threshold is computed.`,
												},
											},
										},
									},

									"conform_action": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"allow"}, false),
										Description:  `Action to take for requests that are under the configured rate limit threshold. Valid option is "allow" only.`,
									},

									"exceed_action": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"redirect", "deny(403)", "deny(404)", "deny(429)", "deny(502)"}, false),
										Description:  `Action to take for requests that are above the configured rate limit threshold, to either deny with a specified HTTP response code, or redirect to a different endpoint. Valid options are "deny()" where valid values for status are 403, 404, 429, and 502, and "redirect" where the redirect parameters come from exceedRedirectOptions below.`,
									},

									"enforce_on_key": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "ALL",
										Description:  `Determines the key to enforce the rateLimitThreshold on`,
										ValidateFunc: validation.StringInSlice([]string{"ALL", "IP", "HTTP_HEADER", "XFF_IP", "HTTP_COOKIE"}, false),
									},

									"enforce_on_key_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Rate limit key name applicable only for the following key types: HTTP_HEADER -- Name of the HTTP header whose value is taken as the key value. HTTP_COOKIE -- Name of the HTTP cookie whose value is taken as the key value.`,
									},

									"ban_threshold": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Can only be specified if the action for the rule is "rate_based_ban". If specified, the key will be banned for the configured 'banDurationSec' when the number of requests that exceed the 'rateLimitThreshold' also exceed this 'banThreshold'.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Number of HTTP(S) requests for calculating the threshold.`,
												},

												"interval_sec": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Interval over which the threshold is computed.`,
												},
											},
										},
									},

									"ban_duration_sec": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Can only be specified if the action for the rule is "rate_based_ban". If specified, determines the time (in seconds) the traffic will continue to be banned by the rate limit after the rate falls below the threshold.`,
									},

									"exceed_redirect_options": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Parameters defining the redirect action that is used as the exceed action. Cannot be specified if the exceed action is not redirect.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:         schema.TypeString,
													Required:     true,
													Description:  `Type of the redirect action.`,
													ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_302", "GOOGLE_RECAPTCHA"}, false),
												},

												"target": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `Target for the redirect action. This is required if the type is EXTERNAL_302 and cannot be specified for GOOGLE_RECAPTCHA.`,
												},
											},
										},
									},
								},
							},
						},

						"redirect_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_302", "GOOGLE_RECAPTCHA"}, false),
										Description:  `Type of the redirect action. Available options: EXTERNAL_302: Must specify the corresponding target field in config. GOOGLE_RECAPTCHA: Cannot specify target field in config.`,
									},

									"target": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Target for the redirect action. This is required if the type is EXTERNAL_302 and cannot be specified for GOOGLE_RECAPTCHA.`,
									},
								},
							},
							Description: `Parameters defining the redirect action. Cannot be specified for any other actions.`,
						},
						//MarkHeader
						"header_action": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Additional actions that are performed on headers.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"request_headers_to_adds": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `The list of request headers to add or overwrite if they're already present.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"header_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The name of the header to set.`,
												},
												"header_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The value to set the named header to.`,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Description: `The set of rules that belong to this policy. There must always be a default rule (rule with priority 2147483647 and match "*"). If no rules are provided when creating a security policy, a default rule with action "allow" will be added.`,
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Fingerprint of this resource.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},

			//Change
			//SetFeature: Cloud Armor for NLB/VMs APIs (item b) -> insert DdosProtectionConfig new object here
			"ddos_protection_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Ddos Protection Config of this security policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ddos_protection": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"STANDARD", "ADVANCED"}, false),
							Description:  `DDOS protection. Supported values include: "STANDARD", "ADVANCED".`,
						},
					},
				},
			},

			"advanced_options_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Advanced Options Config of this security policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"json_parsing": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"DISABLED", "STANDARD"}, false),
							Description:  `JSON body parsing. Supported values include: "DISABLED", "STANDARD".`,
						},
						"json_custom_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `Custom configuration to apply the JSON parsing. Only applicable when JSON parsing is set to STANDARD.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_types": {
										Type:        schema.TypeSet,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `A list of custom Content-Type header values to apply the JSON parsing.`,
									},
								},
							},
						},
						"log_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"NORMAL", "VERBOSE"}, false),
							Description:  `Logging level. Supported values include: "NORMAL", "VERBOSE".`,
						},
					},
				},
			},

			"adaptive_protection_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Adaptive Protection Config of this security policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"layer_7_ddos_defense_config": {
							Type:        schema.TypeList,
							Description: `Layer 7 DDoS Defense Config of this security policy`,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `If set to true, enables CAAP for L7 DDoS detection.`,
									},
									"rule_visibility": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "STANDARD",
										ValidateFunc: validation.StringInSlice([]string{"STANDARD", "PREMIUM"}, false),
										Description:  `Rule visibility. Supported values include: "STANDARD", "PREMIUM".`,
									},
								},
							},
						},
					},
				},
			},
			"recaptcha_options_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `reCAPTCHA configuration options to be applied for the security policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_site_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `A field to supply a reCAPTCHA site key to be used for all the rules using the redirect action with the type of GOOGLE_RECAPTCHA under the security policy. The specified site key needs to be created from the reCAPTCHA API. The user is responsible for the validity of the specified site key. If not specified, a Google-managed site key is used.`,
						},
					},
				},
			},
		},

		UseJSONNumber: true,
	}
}

func resourceComputeRegionSecurityPoliciesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)
	securityPolicy := &compute.SecurityPolicy{
		Name:        sp,
		Description: d.Get("description").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		securityPolicy.Type = v.(string)
	}

	if v, ok := d.GetOk("rule"); ok {
		securityPolicy.Rules = expandSecurityPolicyRules(v.(*schema.Set).List())
	}

	//Change
	if v, ok := d.GetOk("ddos_protection_config"); ok {
		securityPolicy.DdosProtectionConfig = expandSecurityPolicyDdosProtectionConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("advanced_options_config"); ok {
		securityPolicy.AdvancedOptionsConfig = expandSecurityPolicyAdvancedOptionsConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("adaptive_protection_config"); ok {
		securityPolicy.AdaptiveProtectionConfig = expandSecurityPolicyAdaptiveProtectionConfig(v.([]interface{}))
	}

	log.Printf("[DEBUG] RegionSecurityPolicies insert request: %#v", securityPolicy)

	if v, ok := d.GetOk("recaptcha_options_config"); ok {
		securityPolicy.RecaptchaOptionsConfig = expandSecurityPolicyRecaptchaOptionsConfig(v.([]interface{}), d)
	}

	client := config.NewComputeClient(userAgent)

	op, err := client.RegionSecurityPolicies.Insert(project, region, securityPolicy).Do()

	if err != nil {
		return errwrap.Wrapf("Error creating RegionSecurityPolicies: {{err}}", err)
	}

	//id, err := replaceVars(d, config, "projects/{{project}}/global/regionSecurityPolicies/{{name}}")
	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/securityPolicies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Creating RegionSecurityPolicies %q", sp), userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeSecurityPolicyRead(d, meta)
}

func resourceComputeRegionSecurityPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	client := config.NewComputeClient(userAgent)

	securityPolicy, err := client.RegionSecurityPolicies.Get(project, region, sp).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("RegionSecurityPoliciesService %q", d.Id()))
	}

	if err := d.Set("name", securityPolicy.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", securityPolicy.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("type", securityPolicy.Type); err != nil {
		return fmt.Errorf("Error setting type: %s", err)
	}
	if err := d.Set("rule", flattenSecurityPolicyRules(securityPolicy.Rules)); err != nil {
		return err
	}
	if err := d.Set("fingerprint", securityPolicy.Fingerprint); err != nil {
		return fmt.Errorf("Error setting fingerprint: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(securityPolicy.SelfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	//Change
	if err := d.Set("ddos_protection_config", flattenSecurityPolicyDdosProtectionConfig(securityPolicy.DdosProtectionConfig)); err != nil {
		return fmt.Errorf("Error setting ddos_protection_config: %s", err)
	}
	if err := d.Set("advanced_options_config", flattenSecurityPolicyAdvancedOptionsConfig(securityPolicy.AdvancedOptionsConfig)); err != nil {
		return fmt.Errorf("Error setting advanced_options_config: %s", err)
	}

	if err := d.Set("adaptive_protection_config", flattenSecurityPolicyAdaptiveProtectionConfig(securityPolicy.AdaptiveProtectionConfig)); err != nil {
		return fmt.Errorf("Error setting adaptive_protection_config: %s", err)
	}

	if err := d.Set("recaptcha_options_config", flattenSecurityPolicyRecaptchaOptionConfig(securityPolicy.RecaptchaOptionsConfig)); err != nil {
		return fmt.Errorf("Error setting recaptcha_options_config: %s", err)
	}

	return nil
}

func resourceComputeRegionSecurityPoliciesUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	securityPolicy := &compute.SecurityPolicy{
		Fingerprint: d.Get("fingerprint").(string),
	}

	if d.HasChange("type") {
		securityPolicy.Type = d.Get("type").(string)
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "Type")
	}

	if d.HasChange("description") {
		securityPolicy.Description = d.Get("description").(string)
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "Description")
	}
	//Change
	if d.HasChange("ddos_protection_config") {
		securityPolicy.DdosProtectionConfig = expandSecurityPolicyDdosProtectionConfig(d.Get("ddos_protection_config").([]interface{}))
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "DdosProtectionConfig", "ddosProtectionConfig.jsonParsing", "ddosProtectionConfig.jsonCustomConfig", "ddosProtectionConfig.logLevel")
	}

	if d.HasChange("advanced_options_config") {
		securityPolicy.AdvancedOptionsConfig = expandSecurityPolicyAdvancedOptionsConfig(d.Get("advanced_options_config").([]interface{}))
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "AdvancedOptionsConfig", "advancedOptionsConfig.jsonParsing", "advancedOptionsConfig.jsonCustomConfig", "advancedOptionsConfig.logLevel")
	}

	if d.HasChange("adaptive_protection_config") {
		securityPolicy.AdaptiveProtectionConfig = expandSecurityPolicyAdaptiveProtectionConfig(d.Get("adaptive_protection_config").([]interface{}))
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "AdaptiveProtectionConfig", "adaptiveProtectionConfig.layer7DdosDefenseConfig.enable", "adaptiveProtectionConfig.layer7DdosDefenseConfig.ruleVisibility")
	}

	if d.HasChange("recaptcha_options_config") {
		securityPolicy.RecaptchaOptionsConfig = expandSecurityPolicyRecaptchaOptionsConfig(d.Get("recaptcha_options_config").([]interface{}), d)
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "RecaptchaOptionsConfig")
	}

	if len(securityPolicy.ForceSendFields) > 0 {
		client := config.NewComputeClient(userAgent)

		op, err := client.RegionSecurityPolicies.Patch(project, region, sp, securityPolicy).Do()

		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("Error updating RegionSecurityPolicy %q: {{err}}", sp), err)
		}

		err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating RegionSecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		oSet := o.(*schema.Set)
		nSet := n.(*schema.Set)

		oPriorities := map[int64]bool{}
		nPriorities := map[int64]bool{}
		for _, rule := range oSet.List() {
			oPriorities[int64(rule.(map[string]interface{})["priority"].(int))] = true
		}

		for _, rule := range nSet.List() {
			priority := int64(rule.(map[string]interface{})["priority"].(int))
			nPriorities[priority] = true
			if !oPriorities[priority] {
				client := config.NewComputeClient(userAgent)

				// If the rule is in new and its priority does not exist in old, then add it.
				op, err := client.SecurityPolicies.AddRule(project, sp, expandSecurityPolicyRule(rule)).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			} else if !oSet.Contains(rule) {
				client := config.NewComputeClient(userAgent)

				// If the rule is in new, and its priority is in old, but its hash is different than the one in old, update it.
				op, err := client.SecurityPolicies.PatchRule(project, sp, expandSecurityPolicyRule(rule)).Priority(priority).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}

		for _, rule := range oSet.List() {
			priority := int64(rule.(map[string]interface{})["priority"].(int))
			if !nPriorities[priority] {
				client := config.NewComputeClient(userAgent)

				// If the rule's priority is in old but not new, remove it.
				op, err := client.SecurityPolicies.RemoveRule(project, sp).Priority(priority).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceComputeSecurityPolicyRead(d, meta)
}

func resourceComputeRegionSecurityPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	client := config.NewComputeClient(userAgent)

	// Delete the SecurityPolicy
	op, err := client.RegionSecurityPolicies.Delete(project, region, d.Get("name").(string)).Do()
	if err != nil {
		return errwrap.Wrapf("Error deleting RegionSecurityPolicies: {{err}}", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting RegionSecurityPolicies", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceComputeRegionSecurityPoliciesImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	/*if err := parseImportId([]string{"projects/(?P<project>[^/]+)/global/regionSecurityPolicies/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}*/
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/securityPolicies/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	//id, err := replaceVars(d, config, "projects/{{project}}/global/regionSecurityPolicies/{{name}}")
	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/securityPolicies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}