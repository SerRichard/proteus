cwlVersion: v1.2
class: Workflow

inputs:
  wf_input:
    type: string
    default: "Hello!"  # Default input string

outputs:
  workflow_output:
    type: File
    outputSource: generate_parameter/hello_param  # Reference to the output from the generate_parameter step

steps:
  generate_parameter:
    run:
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
    requirements:
      - class: DockerRequirement
        dockerPull: docker/whalesay:latest
        dockerLoad: true
    in:
      input_message: wf_input  # Pass the input from wf_input to the input_message
    out: [hello_param]  # Output file parameter

  consume_parameter:
    run:
      cwlVersion: v1.2
      class: CommandLineTool
      baseCommand: [cowsay]
      inputs:
        message:
          type: string
          inputBinding:
            position: 1
      outputs: []
      arguments: ["$(inputs.message)"]
    requirements:
      - class: DockerRequirement
        dockerPull: docker/whalesay:latest
        dockerLoad: true
    in:
      message: generate_parameter/hello_param  # Use the output from the generate_parameter step as the input
    out: []  # No output declared for this step
