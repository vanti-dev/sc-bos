/**
 * @fileoverview gRPC-Web generated client stub for vanti.bsp.ew.tenants
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')
const proto = {};
proto.vanti = {};
proto.vanti.bsp = {};
proto.vanti.bsp.ew = {};
proto.vanti.bsp.ew.tenants = require('./tenants_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.tenants.TenantApiClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.ListTenantsRequest,
 *   !proto.vanti.bsp.ew.tenants.ListTenantsResponse>}
 */
const methodDescriptor_TenantApi_ListTenants = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/ListTenants',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.ListTenantsRequest,
    proto.vanti.bsp.ew.tenants.ListTenantsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.ListTenantsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.ListTenantsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.ListTenantsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.ListTenantsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.ListTenantsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.listTenants =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/ListTenants',
          request,
          metadata || {},
          methodDescriptor_TenantApi_ListTenants,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.ListTenantsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.ListTenantsResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.listTenants =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/ListTenants',
          request,
          metadata || {},
          methodDescriptor_TenantApi_ListTenants);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.PullTenantsRequest,
 *   !proto.vanti.bsp.ew.tenants.PullTenantsResponse>}
 */
const methodDescriptor_TenantApi_PullTenants = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/PullTenants',
    grpc.web.MethodType.SERVER_STREAMING,
    proto.vanti.bsp.ew.tenants.PullTenantsRequest,
    proto.vanti.bsp.ew.tenants.PullTenantsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.PullTenantsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.PullTenantsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullTenantsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullTenantsResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.pullTenants =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullTenants',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullTenants);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullTenantsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullTenantsResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.pullTenants =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullTenants',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullTenants);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.CreateTenantRequest,
 *   !proto.vanti.bsp.ew.tenants.Tenant>}
 */
const methodDescriptor_TenantApi_CreateTenant = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/CreateTenant',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.CreateTenantRequest,
    proto.vanti.bsp.ew.tenants.Tenant,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.CreateTenantRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Tenant.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.CreateTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Tenant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Tenant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.createTenant =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/CreateTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_CreateTenant,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.CreateTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Tenant>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.createTenant =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/CreateTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_CreateTenant);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.GetTenantRequest,
 *   !proto.vanti.bsp.ew.tenants.Tenant>}
 */
const methodDescriptor_TenantApi_GetTenant = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/GetTenant',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.GetTenantRequest,
    proto.vanti.bsp.ew.tenants.Tenant,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.GetTenantRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Tenant.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.GetTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Tenant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Tenant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.getTenant =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/GetTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_GetTenant,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.GetTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Tenant>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.getTenant =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/GetTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_GetTenant);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.UpdateTenantRequest,
 *   !proto.vanti.bsp.ew.tenants.Tenant>}
 */
const methodDescriptor_TenantApi_UpdateTenant = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/UpdateTenant',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.UpdateTenantRequest,
    proto.vanti.bsp.ew.tenants.Tenant,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.UpdateTenantRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Tenant.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.UpdateTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Tenant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Tenant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.updateTenant =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/UpdateTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_UpdateTenant,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.UpdateTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Tenant>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.updateTenant =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/UpdateTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_UpdateTenant);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.DeleteTenantRequest,
 *   !proto.vanti.bsp.ew.tenants.DeleteTenantResponse>}
 */
const methodDescriptor_TenantApi_DeleteTenant = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/DeleteTenant',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.DeleteTenantRequest,
    proto.vanti.bsp.ew.tenants.DeleteTenantResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.DeleteTenantRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.DeleteTenantResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.DeleteTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.DeleteTenantResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.DeleteTenantResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.deleteTenant =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/DeleteTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_DeleteTenant,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.DeleteTenantRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.DeleteTenantResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.deleteTenant =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/DeleteTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_DeleteTenant);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.PullTenantRequest,
 *   !proto.vanti.bsp.ew.tenants.PullTenantResponse>}
 */
const methodDescriptor_TenantApi_PullTenant = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/PullTenant',
    grpc.web.MethodType.SERVER_STREAMING,
    proto.vanti.bsp.ew.tenants.PullTenantRequest,
    proto.vanti.bsp.ew.tenants.PullTenantResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.PullTenantRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.PullTenantResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullTenantRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullTenantResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.pullTenant =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullTenant);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullTenantRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullTenantResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.pullTenant =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullTenant',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullTenant);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.AddTenantZonesRequest,
 *   !proto.vanti.bsp.ew.tenants.Tenant>}
 */
const methodDescriptor_TenantApi_AddTenantZones = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/AddTenantZones',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.AddTenantZonesRequest,
    proto.vanti.bsp.ew.tenants.Tenant,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.AddTenantZonesRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Tenant.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.AddTenantZonesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Tenant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Tenant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.addTenantZones =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/AddTenantZones',
          request,
          metadata || {},
          methodDescriptor_TenantApi_AddTenantZones,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.AddTenantZonesRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Tenant>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.addTenantZones =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/AddTenantZones',
          request,
          metadata || {},
          methodDescriptor_TenantApi_AddTenantZones);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.RemoveTenantZonesRequest,
 *   !proto.vanti.bsp.ew.tenants.Tenant>}
 */
const methodDescriptor_TenantApi_RemoveTenantZones = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/RemoveTenantZones',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.RemoveTenantZonesRequest,
    proto.vanti.bsp.ew.tenants.Tenant,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.RemoveTenantZonesRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Tenant.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.RemoveTenantZonesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Tenant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Tenant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.removeTenantZones =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/RemoveTenantZones',
          request,
          metadata || {},
          methodDescriptor_TenantApi_RemoveTenantZones,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.RemoveTenantZonesRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Tenant>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.removeTenantZones =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/RemoveTenantZones',
          request,
          metadata || {},
          methodDescriptor_TenantApi_RemoveTenantZones);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.ListSecretsRequest,
 *   !proto.vanti.bsp.ew.tenants.ListSecretsResponse>}
 */
const methodDescriptor_TenantApi_ListSecrets = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/ListSecrets',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.ListSecretsRequest,
    proto.vanti.bsp.ew.tenants.ListSecretsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.ListSecretsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.ListSecretsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.ListSecretsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.ListSecretsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.ListSecretsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.listSecrets =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/ListSecrets',
          request,
          metadata || {},
          methodDescriptor_TenantApi_ListSecrets,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.ListSecretsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.ListSecretsResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.listSecrets =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/ListSecrets',
          request,
          metadata || {},
          methodDescriptor_TenantApi_ListSecrets);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.PullSecretsRequest,
 *   !proto.vanti.bsp.ew.tenants.PullSecretsResponse>}
 */
const methodDescriptor_TenantApi_PullSecrets = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/PullSecrets',
    grpc.web.MethodType.SERVER_STREAMING,
    proto.vanti.bsp.ew.tenants.PullSecretsRequest,
    proto.vanti.bsp.ew.tenants.PullSecretsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.PullSecretsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.PullSecretsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullSecretsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullSecretsResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.pullSecrets =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullSecrets',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullSecrets);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullSecretsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullSecretsResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.pullSecrets =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullSecrets',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullSecrets);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.CreateSecretRequest,
 *   !proto.vanti.bsp.ew.tenants.Secret>}
 */
const methodDescriptor_TenantApi_CreateSecret = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/CreateSecret',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.CreateSecretRequest,
    proto.vanti.bsp.ew.tenants.Secret,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.CreateSecretRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Secret.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.CreateSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Secret)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Secret>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.createSecret =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/CreateSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_CreateSecret,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.CreateSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Secret>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.createSecret =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/CreateSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_CreateSecret);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.VerifySecretRequest,
 *   !proto.vanti.bsp.ew.tenants.Secret>}
 */
const methodDescriptor_TenantApi_VerifySecret = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/VerifySecret',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.VerifySecretRequest,
    proto.vanti.bsp.ew.tenants.Secret,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.VerifySecretRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Secret.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.VerifySecretRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Secret)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Secret>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.verifySecret =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/VerifySecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_VerifySecret,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.VerifySecretRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Secret>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.verifySecret =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/VerifySecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_VerifySecret);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.GetSecretRequest,
 *   !proto.vanti.bsp.ew.tenants.Secret>}
 */
const methodDescriptor_TenantApi_GetSecret = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/GetSecret',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.GetSecretRequest,
    proto.vanti.bsp.ew.tenants.Secret,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.GetSecretRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Secret.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.GetSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Secret)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Secret>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.getSecret =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/GetSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_GetSecret,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.GetSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Secret>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.getSecret =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/GetSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_GetSecret);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.UpdateSecretRequest,
 *   !proto.vanti.bsp.ew.tenants.Secret>}
 */
const methodDescriptor_TenantApi_UpdateSecret = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/UpdateSecret',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.UpdateSecretRequest,
    proto.vanti.bsp.ew.tenants.Secret,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.UpdateSecretRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Secret.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.UpdateSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Secret)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Secret>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.updateSecret =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/UpdateSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_UpdateSecret,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.UpdateSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Secret>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.updateSecret =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/UpdateSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_UpdateSecret);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.DeleteSecretRequest,
 *   !proto.vanti.bsp.ew.tenants.DeleteSecretResponse>}
 */
const methodDescriptor_TenantApi_DeleteSecret = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/DeleteSecret',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.DeleteSecretRequest,
    proto.vanti.bsp.ew.tenants.DeleteSecretResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.DeleteSecretRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.DeleteSecretResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.DeleteSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.DeleteSecretResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.DeleteSecretResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.deleteSecret =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/DeleteSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_DeleteSecret,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.DeleteSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.DeleteSecretResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.deleteSecret =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/DeleteSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_DeleteSecret);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.PullSecretRequest,
 *   !proto.vanti.bsp.ew.tenants.PullSecretResponse>}
 */
const methodDescriptor_TenantApi_PullSecret = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/PullSecret',
    grpc.web.MethodType.SERVER_STREAMING,
    proto.vanti.bsp.ew.tenants.PullSecretRequest,
    proto.vanti.bsp.ew.tenants.PullSecretResponse,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.PullSecretRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.PullSecretResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullSecretRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullSecretResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.pullSecret =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullSecret);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.PullSecretRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.PullSecretResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.pullSecret =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/PullSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_PullSecret);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.tenants.RegenerateSecretRequest,
 *   !proto.vanti.bsp.ew.tenants.Secret>}
 */
const methodDescriptor_TenantApi_RegenerateSecret = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.tenants.TenantApi/RegenerateSecret',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.tenants.RegenerateSecretRequest,
    proto.vanti.bsp.ew.tenants.Secret,
    /**
     * @param {!proto.vanti.bsp.ew.tenants.RegenerateSecretRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.tenants.Secret.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.tenants.RegenerateSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.tenants.Secret)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.tenants.Secret>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.tenants.TenantApiClient.prototype.regenerateSecret =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/RegenerateSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_RegenerateSecret,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.tenants.RegenerateSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.tenants.Secret>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.tenants.TenantApiPromiseClient.prototype.regenerateSecret =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.tenants.TenantApi/RegenerateSecret',
          request,
          metadata || {},
          methodDescriptor_TenantApi_RegenerateSecret);
    };


module.exports = proto.vanti.bsp.ew.tenants;

