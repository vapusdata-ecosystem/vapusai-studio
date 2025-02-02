# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: protos/vapusai-studio/v1alpha1/aimodels.proto
# Protobuf Python Version: 5.29.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    2,
    '',
    'protos/vapusai-studio/v1alpha1/aimodels.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from protos.models.v1alpha1 import vapusai_pb2 as protos_dot_models_dot_v1alpha1_dot_vapusai__pb2
from protos.models.v1alpha1 import common_pb2 as protos_dot_models_dot_v1alpha1_dot_common__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n-protos/vapusai-studio/v1alpha1/aimodels.proto\x12\x10vapusai.v1alpha1\x1a\x1cgoogle/api/annotations.proto\x1a$protos/models/v1alpha1/vapusai.proto\x1a#protos/models/v1alpha1/common.proto\"\x9c\x01\n\x1e\x41IModelNodeConfiguratorRequest\x12H\n\x06\x61\x63tion\x18\x01 \x01(\x0e\x32\x30.vapusai.v1alpha1.AIModelNodeConfiguratorActionsR\x06\x61\x63tion\x12\x30\n\x04spec\x18\x02 \x03(\x0b\x32\x1c.models.v1alpha1.AIModelNodeR\x04spec\"\xc8\x02\n\x13\x41IModelNodeResponse\x12\x34\n\x07\x64m_resp\x18\x01 \x01(\x0b\x32\x1b.models.v1alpha1.DMResponseR\x06\x64mResp\x12Q\n\x06output\x18\x02 \x01(\x0b\x32\x39.vapusai.v1alpha1.AIModelNodeResponse.AIModelNodeResponseR\x06output\x1a\xa7\x01\n\x13\x41IModelNodeResponse\x12\x42\n\x0e\x61i_model_nodes\x18\x01 \x03(\x0b\x32\x1c.models.v1alpha1.AIModelNodeR\x0c\x61iModelNodes\x12L\n\x15\x61vailable_models_pool\x18\x02 \x03(\x0b\x32\x18.models.v1alpha1.MapListR\x13\x61vailableModelsPool\"\x84\x01\n\x18\x41IModelNodeGetterRequest\x12\'\n\x10\x61i_model_node_id\x18\x01 \x01(\tR\raiModelNodeId\x12?\n\x0csearch_param\x18\x02 \x01(\x0b\x32\x1c.models.v1alpha1.SearchParamR\x0bsearchParam*\xa6\x01\n\x1e\x41IModelNodeConfiguratorActions\x12 \n\x1cINVALID_AINODE_CONFIG_ACTION\x10\x00\x12\x1b\n\x17\x43ONFIGURE_AIMODEL_NODES\x10\x01\x12\x17\n\x13PATCH_AIMODEL_NODES\x10\x02\x12\x18\n\x14\x44\x45LETE_AIMODEL_NODES\x10\x03\x12\x12\n\x0eSYNC_AI_MODELS\x10\x04*\xbb\x01\n\x11\x41IModelNodeAction\x12\x1a\n\x16INVALID_AIAGENT_ACTION\x10\x00\x12\x14\n\x10GENERATE_CONTENT\x10\x01\x12\x17\n\x13GENERATE_TRANSCRIPT\x10\x02\x12\x18\n\x14GENERATE_DATAQUERIES\x10\x03\x12\x13\n\x0fGENERATE_IMAGES\x10\x04\x12\x13\n\x0fGENERATE_VIDEOS\x10\x05\x12\x17\n\x13GENERATE_EMBEDDINGS\x10\x06\x32\xab\x02\n\x08\x41IModels\x12\x89\x01\n\x07Manager\x12\x30.vapusai.v1alpha1.AIModelNodeConfiguratorRequest\x1a%.vapusai.v1alpha1.AIModelNodeResponse\"%\x82\xd3\xe4\x93\x02\x1f\"\x1a/api/v1alpha1/models-nodes:\x01*\x12\x92\x01\n\x06Getter\x12*.vapusai.v1alpha1.AIModelNodeGetterRequest\x1a%.vapusai.v1alpha1.AIModelNodeResponse\"5\x82\xd3\xe4\x93\x02/\x12-/api/v1alpha1/models-nodes/{ai_model_node_id}B\xbd\x01\n\x14\x63om.vapusai.v1alpha1B\rAimodelsProtoP\x01Z5github.com/vapusdata-oss/apis/protos/vapusai/v1alpha1\xa2\x02\x03VXX\xaa\x02\x10Vapusai.V1alpha1\xca\x02\x10Vapusai\\V1alpha1\xe2\x02\x1cVapusai\\V1alpha1\\GPBMetadata\xea\x02\x11Vapusai::V1alpha1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'protos.vapusai_studio.v1alpha1.aimodels_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\024com.vapusai.v1alpha1B\rAimodelsProtoP\001Z5github.com/vapusdata-oss/apis/protos/vapusai/v1alpha1\242\002\003VXX\252\002\020Vapusai.V1alpha1\312\002\020Vapusai\\V1alpha1\342\002\034Vapusai\\V1alpha1\\GPBMetadata\352\002\021Vapusai::V1alpha1'
  _globals['_AIMODELS'].methods_by_name['Manager']._loaded_options = None
  _globals['_AIMODELS'].methods_by_name['Manager']._serialized_options = b'\202\323\344\223\002\037\"\032/api/v1alpha1/models-nodes:\001*'
  _globals['_AIMODELS'].methods_by_name['Getter']._loaded_options = None
  _globals['_AIMODELS'].methods_by_name['Getter']._serialized_options = b'\202\323\344\223\002/\022-/api/v1alpha1/models-nodes/{ai_model_node_id}'
  _globals['_AIMODELNODECONFIGURATORACTIONS']._serialized_start=798
  _globals['_AIMODELNODECONFIGURATORACTIONS']._serialized_end=964
  _globals['_AIMODELNODEACTION']._serialized_start=967
  _globals['_AIMODELNODEACTION']._serialized_end=1154
  _globals['_AIMODELNODECONFIGURATORREQUEST']._serialized_start=173
  _globals['_AIMODELNODECONFIGURATORREQUEST']._serialized_end=329
  _globals['_AIMODELNODERESPONSE']._serialized_start=332
  _globals['_AIMODELNODERESPONSE']._serialized_end=660
  _globals['_AIMODELNODERESPONSE_AIMODELNODERESPONSE']._serialized_start=493
  _globals['_AIMODELNODERESPONSE_AIMODELNODERESPONSE']._serialized_end=660
  _globals['_AIMODELNODEGETTERREQUEST']._serialized_start=663
  _globals['_AIMODELNODEGETTERREQUEST']._serialized_end=795
  _globals['_AIMODELS']._serialized_start=1157
  _globals['_AIMODELS']._serialized_end=1456
# @@protoc_insertion_point(module_scope)
