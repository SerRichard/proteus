cwlVersion: v1.2
class: CommandLineTool
id: param-ref
requirements: 
  - class: DockerRequirement 
    dockerPull: ubuntu:20.04
    dockerOutputDirectory: /tmp

  - class: ResourceRequirement 
    outdirMin: 1Gi
baseCommand: [tar, --extract]
inputs:
  tarfile:
    type: File
    inputBinding:
      prefix: --file
outputs:
  extracted_file:
    type: File
    outputBinding:
      glob: $(inputs.extractfile)
