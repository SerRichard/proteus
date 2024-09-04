cwlVersion: v1.0
class: CommandLineTool
id: echo-tool
requirements: 
  - class: DockerRequirement 
    dockerPull: ubuntu:20.04
    dockerOutputDirectory: /tmp

  - class: ResourceRequirement 
    outdirMin: 1Gi

baseCommand: echo
inputs:
  message:
    type: string
    inputBinding:
      position: 1
    default: "Hello!"  # Default input string

outputs: 
  hello_param:
    type: File
    outputBinding:
      glob: /tmp/hello_world.txt