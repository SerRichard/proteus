cwlVersion: v1.2
class: Workflow
inputs:
  wf_input:
    type: string
    default: "Hello!"
outputs:
  workflow_output:
    type: File
    outputSource: generate_parameter/hello_param
steps:
  generate_parameter:
    run: whalesay.cwl
    requirements:
      - class: DockerRequirement
        dockerPull: docker/whalesay:latest
        dockerLoad: true
    in:
      input_message: wf_input
    out: [hello_param]
  
  consume_parameter:
    run: print-message.cwl
    in:
      message: generate_parameter/hello_param
    requirements:
      - class: DockerRequirement
        dockerPull: docker/whalesay:latest
        dockerLoad: true
    out: []


