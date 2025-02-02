# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: protos/models/v1alpha1/users.proto
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
    'protos/models/v1alpha1/users.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from protos.models.v1alpha1 import common_pb2 as protos_dot_models_dot_v1alpha1_dot_common__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\"protos/models/v1alpha1/users.proto\x12\x0fmodels.v1alpha1\x1a#protos/models/v1alpha1/common.proto\"\xd4\x04\n\x04User\x12!\n\x0c\x64isplay_name\x18\x01 \x01(\tR\x0b\x64isplayName\x12\x17\n\x07user_id\x18\x02 \x01(\tR\x06userId\x12\x14\n\x05\x65mail\x18\x03 \x01(\tR\x05\x65mail\x12T\n\x12organization_roles\x18\x04 \x03(\x0b\x32%.models.v1alpha1.UserOrganizationRoleR\x11organizationRoles\x12\x1b\n\tinvite_id\x18\x05 \x01(\tR\x08inviteId\x12\x43\n\x0cstudio_roles\x18\x06 \x03(\x0e\x32 .models.v1alpha1.StudioUserRolesR\x0bstudioRoles\x12\x16\n\x06status\x18\x07 \x01(\tR\x06status\x12\x1d\n\ninvited_on\x18\x08 \x01(\x03R\tinvitedOn\x12*\n\x11invite_expires_on\x18\t \x01(\x03R\x0finviteExpiresOn\x12\x1d\n\nfirst_name\x18\n \x01(\tR\tfirstName\x12\x1b\n\tlast_name\x18\x0b \x01(\tR\x08lastName\x12\x42\n\x0cinvited_type\x18\x0c \x01(\x0e\x32\x1f.models.v1alpha1.UserInviteTypeR\x0binvitedType\x12\'\n\x0fstudio_policies\x18\r \x03(\tR\x0estudioPolicies\x12\x36\n\x07profile\x18\x0e \x01(\x0b\x32\x1c.models.v1alpha1.UserProfileR\x07profile\"\x7f\n\x0bUserProfile\x12\x36\n\taddresses\x18\x01 \x03(\x0b\x32\x18.models.v1alpha1.AddressR\taddresses\x12\x16\n\x06\x61vatar\x18\x02 \x01(\tR\x06\x61vatar\x12 \n\x0b\x64\x65scription\x18\x03 \x01(\tR\x0b\x64\x65scription\"8\n\tUserShort\x12\x17\n\x07user_id\x18\x01 \x01(\tR\x06userId\x12\x12\n\x04role\x18\x02 \x01(\tR\x04role\"V\n\nUserValidT\x12)\n\x04user\x18\x01 \x01(\x0b\x32\x15.models.v1alpha1.UserR\x04user\x12\x1d\n\nvalid_till\x18\x02 \x01(\tR\tvalidTill\"\xa6\x01\n\x13OrganizationUserOps\x12\x17\n\x07user_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n\nvalid_till\x18\x02 \x01(\x03R\tvalidTill\x12\x12\n\x04role\x18\x03 \x03(\tR\x04role\x12\x14\n\x05\x61\x64\x64\x65\x64\x18\x04 \x01(\x08R\x05\x61\x64\x64\x65\x64\x12-\n\x13invite_if_not_found\x18\x05 \x01(\x08R\x10inviteIfNotFound\"\x88\x02\n\nUserInvite\x12\x14\n\x05\x65mail\x18\x01 \x01(\tR\x05\x65mail\x12\x1d\n\ninvited_on\x18\x02 \x01(\x03R\tinvitedOn\x12*\n\x11invite_expires_on\x18\x03 \x01(\x03R\x0finviteExpiresOn\x12\x1b\n\tinvite_id\x18\x04 \x01(\tR\x08inviteId\x12\x12\n\x04role\x18\x05 \x01(\tR\x04role\x12T\n\x12organization_roles\x18\x06 \x03(\x0b\x32%.models.v1alpha1.UserOrganizationRoleR\x11organizationRoles\x12\x12\n\x04name\x18\x07 \x01(\tR\x04name\"\x82\x01\n\x04Team\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12+\n\x05users\x18\x02 \x03(\x0b\x32\x15.models.v1alpha1.UserR\x05users\x12 \n\x0b\x64\x65scription\x18\x03 \x01(\tR\x0b\x64\x65scription\x12\x17\n\x07team_id\x18\x04 \x01(\tR\x06teamId\"\xcc\x01\n\x14UserOrganizationRole\x12\'\n\x0forganization_id\x18\x01 \x01(\tR\x0eorganizationId\x12\x12\n\x04role\x18\x02 \x03(\tR\x04role\x12\x1d\n\ninvited_on\x18\x03 \x01(\x03R\tinvitedOn\x12\x1d\n\nvalid_till\x18\x04 \x01(\x03R\tvalidTill\x12\x1a\n\x08policies\x18\x05 \x03(\tR\x08policies\x12\x1d\n\nis_default\x18\x06 \x01(\x08R\tisDefault\"\xa3\x01\n\nVapusRoles\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\x16\n\x06\x61\x63tion\x18\x02 \x01(\tR\x06\x61\x63tion\x12\x10\n\x03\x61rn\x18\x03 \x01(\tR\x03\x61rn\x12\x35\n\x05scope\x18\x04 \x01(\x0e\x32\x1f.models.v1alpha1.VapusRoleScopeR\x05scope\x12 \n\x0b\x64\x65scription\x18\x05 \x01(\tR\x0b\x64\x65scription\"\x99\x01\n\x0cRefreshToken\x12\x15\n\x06jwt_id\x18\x01 \x01(\tR\x05jwtId\x12\x1d\n\nvalid_till\x18\x02 \x01(\tR\tvalidTill\x12\x16\n\x06status\x18\x03 \x01(\tR\x06status\x12\x17\n\x07user_id\x18\x04 \x01(\tR\x06userId\x12\"\n\x0corganization\x18\x05 \x01(\tR\x0corganization\"\xac\x02\n\x10VapusResourceArn\x12#\n\rresource_name\x18\x01 \x01(\tR\x0cresourceName\x12\x1f\n\x0bresource_id\x18\x02 \x01(\tR\nresourceId\x12!\n\x0cresource_arn\x18\x03 \x01(\tR\x0bresourceArn\x12\x44\n\rallowed_rules\x18\x04 \x03(\x0b\x32\x1f.models.v1alpha1.ResourceAclMapR\x0c\x61llowedRules\x12\x44\n\rblocked_rules\x18\x05 \x03(\x0b\x32\x1f.models.v1alpha1.ResourceAclMapR\x0c\x62lockedRules\x12#\n\rblocked_users\x18\x06 \x03(\tR\x0c\x62lockedUsers\"J\n\x0eResourceAclMap\x12\"\n\x0corganization\x18\x01 \x01(\tR\x0corganization\x12\x14\n\x05users\x18\x02 \x03(\tR\x05users*z\n\x0fStudioUserRoles\x12\r\n\tANONYMOUS\x10\x00\x12\x11\n\rSTUDIO_OWNERS\x10\x01\x12\x14\n\x10STUDIO_OPERATORS\x10\x02\x12\x0e\n\nORG_OWNERS\x10\x03\x12\r\n\tORG_USERS\x10\x04\x12\x10\n\x0cSTUDIO_USERS\x10\x05*K\n\x0eVapusRoleScope\x12\x10\n\x0c\x44OMAIN_ROLES\x10\x00\x12\x10\n\x0cSTUDIO_ROLES\x10\x01\x12\x15\n\x11MARKETPLACE_ROLES\x10\x02*I\n\x0eUserInviteType\x12\x12\n\x0eREQUEST_ACCESS\x10\x00\x12\x11\n\rINVITE_ACCESS\x10\x01\x12\x10\n\x0cSTUDIO_SETUP\x10\x02\x42\xb4\x01\n\x13\x63om.models.v1alpha1B\nUsersProtoP\x01Z4github.com/vapusdata-oss/apis/protos/models/v1alpha1\xa2\x02\x03MXX\xaa\x02\x0fModels.V1alpha1\xca\x02\x0fModels\\V1alpha1\xe2\x02\x1bModels\\V1alpha1\\GPBMetadata\xea\x02\x10Models::V1alpha1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'protos.models.v1alpha1.users_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\023com.models.v1alpha1B\nUsersProtoP\001Z4github.com/vapusdata-oss/apis/protos/models/v1alpha1\242\002\003MXX\252\002\017Models.V1alpha1\312\002\017Models\\V1alpha1\342\002\033Models\\V1alpha1\\GPBMetadata\352\002\020Models::V1alpha1'
  _globals['_STUDIOUSERROLES']._serialized_start=2443
  _globals['_STUDIOUSERROLES']._serialized_end=2565
  _globals['_VAPUSROLESCOPE']._serialized_start=2567
  _globals['_VAPUSROLESCOPE']._serialized_end=2642
  _globals['_USERINVITETYPE']._serialized_start=2644
  _globals['_USERINVITETYPE']._serialized_end=2717
  _globals['_USER']._serialized_start=93
  _globals['_USER']._serialized_end=689
  _globals['_USERPROFILE']._serialized_start=691
  _globals['_USERPROFILE']._serialized_end=818
  _globals['_USERSHORT']._serialized_start=820
  _globals['_USERSHORT']._serialized_end=876
  _globals['_USERVALIDT']._serialized_start=878
  _globals['_USERVALIDT']._serialized_end=964
  _globals['_ORGANIZATIONUSEROPS']._serialized_start=967
  _globals['_ORGANIZATIONUSEROPS']._serialized_end=1133
  _globals['_USERINVITE']._serialized_start=1136
  _globals['_USERINVITE']._serialized_end=1400
  _globals['_TEAM']._serialized_start=1403
  _globals['_TEAM']._serialized_end=1533
  _globals['_USERORGANIZATIONROLE']._serialized_start=1536
  _globals['_USERORGANIZATIONROLE']._serialized_end=1740
  _globals['_VAPUSROLES']._serialized_start=1743
  _globals['_VAPUSROLES']._serialized_end=1906
  _globals['_REFRESHTOKEN']._serialized_start=1909
  _globals['_REFRESHTOKEN']._serialized_end=2062
  _globals['_VAPUSRESOURCEARN']._serialized_start=2065
  _globals['_VAPUSRESOURCEARN']._serialized_end=2365
  _globals['_RESOURCEACLMAP']._serialized_start=2367
  _globals['_RESOURCEACLMAP']._serialized_end=2441
# @@protoc_insertion_point(module_scope)
