cwlVersion: v1.2
class: CommandLineTool
baseCommand: [sh, -c]
requirements:
  - class: DockerRequirement
    dockerPull: docker/whalesay:latest
    dockerLoad: true
inputs:
  input_message:
    type: string
    inputBinding:
      position: 1
outputs:
  hello_param:
    type: File
    outputBinding:
      glob: /tmp/hello_world.txt
arguments: ["sleep 1; echo -n $(inputs.input_message) > /tmp/hello_world.txt"]
