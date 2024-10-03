cwlVersion: v1.2
class: CommandLineTool
id: array-inputs
requirements: 
  - class: DockerRequirement 
    dockerPull: ubuntu:20.04
    dockerOutputDirectory: /tmp

  - class: ResourceRequirement 
    outdirMin: 1Gi

inputs:
  filesA:
    type: string[]
    inputBinding:
      prefix: -A
      position: 1

  filesB:
    type: array
    inputBinding:
      position: 2

  filesC:
    type: string[]
    inputBinding:
      prefix: -C=
      itemSeparator: ","
      separate: false
      position: 4

outputs:
  example_out:
    type: File
    outputSource: tmp/stdout_data

baseCommand: echo