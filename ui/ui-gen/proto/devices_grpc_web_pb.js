/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v5.29.3
// source: devices.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var traits_metadata_pb = require('@smart-core-os/sc-api-grpc-web/traits/metadata_pb.js')

var types_change_pb = require('@smart-core-os/sc-api-grpc-web/types/change_pb.js')

var types_time_period_pb = require('@smart-core-os/sc-api-grpc-web/types/time/period_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./devices_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.DevicesApiClient =
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
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.DevicesApiPromiseClient =
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
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.ListDevicesRequest,
 *   !proto.smartcore.bos.ListDevicesResponse>}
 */
const methodDescriptor_DevicesApi_ListDevices = new grpc.web.MethodDescriptor(
  '/smartcore.bos.DevicesApi/ListDevices',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ListDevicesRequest,
  proto.smartcore.bos.ListDevicesResponse,
  /**
   * @param {!proto.smartcore.bos.ListDevicesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListDevicesResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ListDevicesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.ListDevicesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ListDevicesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiClient.prototype.listDevices =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/ListDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_ListDevices,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ListDevicesRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ListDevicesResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.DevicesApiPromiseClient.prototype.listDevices =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/ListDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_ListDevices);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullDevicesRequest,
 *   !proto.smartcore.bos.PullDevicesResponse>}
 */
const methodDescriptor_DevicesApi_PullDevices = new grpc.web.MethodDescriptor(
  '/smartcore.bos.DevicesApi/PullDevices',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullDevicesRequest,
  proto.smartcore.bos.PullDevicesResponse,
  /**
   * @param {!proto.smartcore.bos.PullDevicesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullDevicesResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullDevicesRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullDevicesResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiClient.prototype.pullDevices =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.DevicesApi/PullDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_PullDevices);
};


/**
 * @param {!proto.smartcore.bos.PullDevicesRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullDevicesResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiPromiseClient.prototype.pullDevices =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.DevicesApi/PullDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_PullDevices);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.GetDevicesMetadataRequest,
 *   !proto.smartcore.bos.DevicesMetadata>}
 */
const methodDescriptor_DevicesApi_GetDevicesMetadata = new grpc.web.MethodDescriptor(
  '/smartcore.bos.DevicesApi/GetDevicesMetadata',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetDevicesMetadataRequest,
  proto.smartcore.bos.DevicesMetadata,
  /**
   * @param {!proto.smartcore.bos.GetDevicesMetadataRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.DevicesMetadata.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetDevicesMetadataRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.DevicesMetadata)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.DevicesMetadata>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiClient.prototype.getDevicesMetadata =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/GetDevicesMetadata',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_GetDevicesMetadata,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetDevicesMetadataRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.DevicesMetadata>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.DevicesApiPromiseClient.prototype.getDevicesMetadata =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/GetDevicesMetadata',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_GetDevicesMetadata);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullDevicesMetadataRequest,
 *   !proto.smartcore.bos.PullDevicesMetadataResponse>}
 */
const methodDescriptor_DevicesApi_PullDevicesMetadata = new grpc.web.MethodDescriptor(
  '/smartcore.bos.DevicesApi/PullDevicesMetadata',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullDevicesMetadataRequest,
  proto.smartcore.bos.PullDevicesMetadataResponse,
  /**
   * @param {!proto.smartcore.bos.PullDevicesMetadataRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullDevicesMetadataResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullDevicesMetadataRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullDevicesMetadataResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiClient.prototype.pullDevicesMetadata =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.DevicesApi/PullDevicesMetadata',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_PullDevicesMetadata);
};


/**
 * @param {!proto.smartcore.bos.PullDevicesMetadataRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullDevicesMetadataResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiPromiseClient.prototype.pullDevicesMetadata =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.DevicesApi/PullDevicesMetadata',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_PullDevicesMetadata);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.GetDownloadDevicesUrlRequest,
 *   !proto.smartcore.bos.DownloadDevicesUrl>}
 */
const methodDescriptor_DevicesApi_GetDownloadDevicesUrl = new grpc.web.MethodDescriptor(
  '/smartcore.bos.DevicesApi/GetDownloadDevicesUrl',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetDownloadDevicesUrlRequest,
  proto.smartcore.bos.DownloadDevicesUrl,
  /**
   * @param {!proto.smartcore.bos.GetDownloadDevicesUrlRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.DownloadDevicesUrl.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetDownloadDevicesUrlRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.DownloadDevicesUrl)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.DownloadDevicesUrl>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiClient.prototype.getDownloadDevicesUrl =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/GetDownloadDevicesUrl',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_GetDownloadDevicesUrl,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetDownloadDevicesUrlRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.DownloadDevicesUrl>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.DevicesApiPromiseClient.prototype.getDownloadDevicesUrl =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/GetDownloadDevicesUrl',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_GetDownloadDevicesUrl);
};


module.exports = proto.smartcore.bos;

