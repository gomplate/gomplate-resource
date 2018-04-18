# gomplate-resource

Concourse CI resource for templating files with [gomplate][]

## Examples

```yaml
resource_types:
- name: gomplate
  type: docker-image
  source:
    repository: gomplate/gomplate-resource
    tag: latest

resources:
- name: gomplate
  type: gomplate

jobs:
- name: deploy-things
  plan:
  - get: things
    trigger: true
  - get: datasources
  - get: rendered-files
    resource: gomplate
    params:
      datasources:
        - data=datasources/data.yml
      inputFiles:
        - things/file1.in
        - things/file2.in
      # these are prefixed by the destination dir, i.e. `rendered-files/`
      outputFiles:
        - file1.out
        - file2.out
  - get: rendered-dir
    resource: gomplate
    params:
      inputDir: things
      # outputDir defaults to the destination directory
```

[gomplate]: https://gomplate.hairyhenderson.ca/
