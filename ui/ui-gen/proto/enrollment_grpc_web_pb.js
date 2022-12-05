/**
 * @fileoverview gRPC-Web generated client stub for vanti.bsp.ew
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.vanti = {};
proto.vanti.bsp = {};
proto.vanti.bsp.ew = require('./enrollment_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.EnrollmentApiClient =
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
proto.vanti.bsp.ew.EnrollmentApiPromiseClient =
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
 *   !proto.vanti.bsp.ew.GetEnrollmentRequest,
 *   !proto.vanti.bsp.ew.Enrollment>}
 */
const methodDescriptor_EnrollmentApi_GetEnrollment = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.EnrollmentApi/GetEnrollment',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.GetEnrollmentRequest,
    proto.vanti.bsp.ew.Enrollment,
    /**
     * @param {!proto.vanti.bsp.ew.GetEnrollmentRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Enrollment.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.GetEnrollmentRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Enrollment)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Enrollment>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.EnrollmentApiClient.prototype.getEnrollment =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.EnrollmentApi/GetEnrollment',
          request,
          metadata || {},
          methodDescriptor_EnrollmentApi_GetEnrollment,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.GetEnrollmentRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Enrollment>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.EnrollmentApiPromiseClient.prototype.getEnrollment =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.EnrollmentApi/GetEnrollment',
          request,
          metadata || {},
          methodDescriptor_EnrollmentApi_GetEnrollment);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.CreateEnrollmentRequest,
 *   !proto.vanti.bsp.ew.Enrollment>}
 */
const methodDescriptor_EnrollmentApi_CreateEnrollment = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.EnrollmentApi/CreateEnrollment',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.CreateEnrollmentRequest,
    proto.vanti.bsp.ew.Enrollment,
    /**
     * @param {!proto.vanti.bsp.ew.CreateEnrollmentRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Enrollment.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.CreateEnrollmentRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Enrollment)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Enrollment>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.EnrollmentApiClient.prototype.createEnrollment =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.EnrollmentApi/CreateEnrollment',
          request,
          metadata || {},
          methodDescriptor_EnrollmentApi_CreateEnrollment,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.CreateEnrollmentRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Enrollment>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.EnrollmentApiPromiseClient.prototype.createEnrollment =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.EnrollmentApi/CreateEnrollment',
          request,
          metadata || {},
          methodDescriptor_EnrollmentApi_CreateEnrollment);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.DeleteEnrollmentRequest,
 *   !proto.vanti.bsp.ew.Enrollment>}
 */
const methodDescriptor_EnrollmentApi_DeleteEnrollment = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.EnrollmentApi/DeleteEnrollment',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.DeleteEnrollmentRequest,
    proto.vanti.bsp.ew.Enrollment,
    /**
     * @param {!proto.vanti.bsp.ew.DeleteEnrollmentRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Enrollment.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.DeleteEnrollmentRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Enrollment)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Enrollment>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.EnrollmentApiClient.prototype.deleteEnrollment =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.EnrollmentApi/DeleteEnrollment',
          request,
          metadata || {},
          methodDescriptor_EnrollmentApi_DeleteEnrollment,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.DeleteEnrollmentRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Enrollment>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.EnrollmentApiPromiseClient.prototype.deleteEnrollment =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.EnrollmentApi/DeleteEnrollment',
          request,
          metadata || {},
          methodDescriptor_EnrollmentApi_DeleteEnrollment);
    };


module.exports = proto.vanti.bsp.ew;

