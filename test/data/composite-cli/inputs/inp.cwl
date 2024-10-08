cwlVersion: v1.2
class: CommandLineTool
baseCommand: echo
inputs:
  example_flag:
    type: boolean
    inputBinding:
      position: 1
      prefix: -f
  example_string:
    type: string
    inputBinding:
      position: 3
      prefix: --example-string
  example_int:
    type: int
    inputBinding:
      position: 2
      prefix: -i
      separate: false
  example_file:
    type: File
    inputBinding:
      prefix: --file=
      separate: false
      position: 4
outputs: []
id: echo-tool
requirements: 
  - class: DockerRequirement 
    dockerPull: ubuntu:20.04
    dockerOutputDirectory: /tmp