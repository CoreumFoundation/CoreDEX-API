// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v1.181.2
//   protoc               v5.29.3
// source: domain/metadata/metadata.proto
/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Timestamp } from "../../google/protobuf/timestamp";
export const protobufPackage = "metadata";
export var Network;
(function (Network) {
    /** NETWORK_DO_NOT_USE - Avoid a default value since a default and a bug are not distinguishable */
    Network[Network["NETWORK_DO_NOT_USE"] = 0] = "NETWORK_DO_NOT_USE";
    Network[Network["MAINNET"] = 1] = "MAINNET";
    Network[Network["TESTNET"] = 2] = "TESTNET";
    Network[Network["DEVNET"] = 3] = "DEVNET";
    Network[Network["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(Network || (Network = {}));
export function networkFromJSON(object) {
    switch (object) {
        case 0:
        case "NETWORK_DO_NOT_USE":
            return Network.NETWORK_DO_NOT_USE;
        case 1:
        case "MAINNET":
            return Network.MAINNET;
        case 2:
        case "TESTNET":
            return Network.TESTNET;
        case 3:
        case "DEVNET":
            return Network.DEVNET;
        case -1:
        case "UNRECOGNIZED":
        default:
            return Network.UNRECOGNIZED;
    }
}
export function networkToJSON(object) {
    switch (object) {
        case Network.NETWORK_DO_NOT_USE:
            return "NETWORK_DO_NOT_USE";
        case Network.MAINNET:
            return "MAINNET";
        case Network.TESTNET:
            return "TESTNET";
        case Network.DEVNET:
            return "DEVNET";
        case Network.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
function createBaseMetaData() {
    return { Network: 0, UpdatedAt: undefined, CreatedAt: undefined };
}
export const MetaData = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.Network !== 0) {
            writer.uint32(8).int32(message.Network);
        }
        if (message.UpdatedAt !== undefined) {
            Timestamp.encode(toTimestamp(message.UpdatedAt), writer.uint32(18).fork()).ldelim();
        }
        if (message.CreatedAt !== undefined) {
            Timestamp.encode(toTimestamp(message.CreatedAt), writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMetaData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    if (tag !== 8) {
                        break;
                    }
                    message.Network = reader.int32();
                    continue;
                case 2:
                    if (tag !== 18) {
                        break;
                    }
                    message.UpdatedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
                    continue;
                case 3:
                    if (tag !== 26) {
                        break;
                    }
                    message.CreatedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
                    continue;
            }
            if ((tag & 7) === 4 || tag === 0) {
                break;
            }
            reader.skipType(tag & 7);
        }
        return message;
    },
    fromJSON(object) {
        return {
            Network: isSet(object.Network) ? networkFromJSON(object.Network) : 0,
            UpdatedAt: isSet(object.UpdatedAt) ? fromJsonTimestamp(object.UpdatedAt) : undefined,
            CreatedAt: isSet(object.CreatedAt) ? fromJsonTimestamp(object.CreatedAt) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        if (message.Network !== 0) {
            obj.Network = networkToJSON(message.Network);
        }
        if (message.UpdatedAt !== undefined) {
            obj.UpdatedAt = message.UpdatedAt.toISOString();
        }
        if (message.CreatedAt !== undefined) {
            obj.CreatedAt = message.CreatedAt.toISOString();
        }
        return obj;
    },
    create(base) {
        return MetaData.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseMetaData();
        message.Network = (_a = object.Network) !== null && _a !== void 0 ? _a : 0;
        message.UpdatedAt = (_b = object.UpdatedAt) !== null && _b !== void 0 ? _b : undefined;
        message.CreatedAt = (_c = object.CreatedAt) !== null && _c !== void 0 ? _c : undefined;
        return message;
    },
};
function toTimestamp(date) {
    const seconds = Math.trunc(date.getTime() / 1000);
    const nanos = (date.getTime() % 1000) * 1000000;
    return { seconds, nanos };
}
function fromTimestamp(t) {
    let millis = (t.seconds || 0) * 1000;
    millis += (t.nanos || 0) / 1000000;
    return new globalThis.Date(millis);
}
function fromJsonTimestamp(o) {
    if (o instanceof globalThis.Date) {
        return o;
    }
    else if (typeof o === "string") {
        return new globalThis.Date(o);
    }
    else {
        return fromTimestamp(Timestamp.fromJSON(o));
    }
}
function isSet(value) {
    return value !== null && value !== undefined;
}
