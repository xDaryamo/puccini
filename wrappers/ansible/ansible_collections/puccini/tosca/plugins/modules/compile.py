# https://docs.ansible.com/ansible/latest/dev_guide/developing_modules_general.html

from ansible.module_utils.basic import AnsibleModule
import puccini.tosca, ard

DOCUMENTATION = r'''
---
module: tosca

short_description: TOSCA

version_added: "1.0.0"

description: Compile TOSCA service template to Clout.

options:
    service_template:
        description: A URL or file path to a TOSCA service template. Can be a CSAR or YAML file.
        required: true
        type: str
    debug:
        description: Set to true to enable the raw "clout" output.
        required: false
        type: bool

author:
    - Puccini (@puccini)
'''

EXAMPLES = r'''
- name: Compile TOSCA service template
  puccini.tosca.compile:
    service_template: ../../examples/tosca/requirements-and-capabilities.yaml
    inputs: {}
  register: service
'''

RETURN = r'''
node_templates: [...]
'''

def run_module():
    module_args = dict(
        service_template=dict(type='str', required=True),
        inputs=dict(type='dict', required=False, default=dict()),
        debug=dict(type='bool', required=False, default=False),
    )

    module = AnsibleModule(
        argument_spec=module_args,
        supports_check_mode=True
    )

    result = dict()

    if module.check_mode:
        module.exit_json(**result)

    service_template = module.params.get('service_template', '')
    inputs = module.params.get('inputs', {})
    try:
        clout = puccini.tosca.compile(service_template, inputs)
    except Exception as e:
        module.fail_json(msg=str(e))

    if module.params['debug']:
        result['clout'] = clout

    result['node_templates'] = []
    for vertex in clout['vertexes'].values():
        try:
            if vertex.get('metadata', {}).get('puccini', {}).get('kind', '') == 'NodeTemplate':
                node_template = vertex.get('properties', {})
                result['node_templates'].append(node_template)
        except:
            pass

    result = ard.cjson.convert_to(result)
    module.exit_json(**result)

if __name__ == '__main__':
    run_module()
