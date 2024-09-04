cwlVersion: v1.2
class: CommandLineTool
baseCommand: [cowsay]
inputs:
  message:
    type: string
    inputBinding:
      position: 1
outputs: []
requirements:
  - class: DockerRequirement
    dockerPull: docker/whalesay:latest
    dockerLoad: true
arguments: ["$(inputs.message)"]