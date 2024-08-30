cwlVersion: v1.0
class: CommandLineTool
id: echo-tool
requirements: 
  - class: DockerRequirement 
    dockerPull: ubuntu:20.04
baseCommand: echo
inputs:
  message:
    type: string
    inputBinding:
      position: 1
    default: "Hello!"  # Default input string
outputs: []