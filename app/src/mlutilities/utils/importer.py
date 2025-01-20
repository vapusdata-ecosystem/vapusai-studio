import os
import sys

def proto_importer():
    BASE_DIR = os.path.dirname(os.path.abspath(__file__))

    # Construct the path to the protobuf directory relative to BASE_DIR
    PROTOBUF_DIR = os.path.join(BASE_DIR, '../../../../apis', 'gen-python')

    # Add the protobuf directory to the system path
    sys.path.append(os.path.abspath(PROTOBUF_DIR))