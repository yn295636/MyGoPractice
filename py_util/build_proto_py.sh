#!/usr/bin/env bash
if [[ ! -d pyvenv ]]; then
  virtualenv -p python3 pyvenv
fi
source pyvenv/bin/activate
pip install --upgrade -r py_util/requirements.txt
python -m grpc_tools.protoc -I=. \
 --python_out=. \
 --grpc_python_out=. \
 proto/sample_service/sample_service.proto
python -m grpc_tools.protoc -I=. \
 --python_out=. \
 --grpc_python_out=. \
 proto/greeter_service/greeter_service.proto