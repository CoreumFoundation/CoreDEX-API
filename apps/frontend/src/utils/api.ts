import axios, { AxiosResponse } from "axios";
import { useStore } from "@/state/store";

export enum APIMethod {
  GET = "get",
  POST = "post",
  UPDATE = "put",
  DELETE = "delete",
  POST_FILE = "postFile",
}

export const request = async (
  body: any,
  endpoint: string,
  method: APIMethod
): Promise<AxiosResponse> => {
  let headers: any = {
    "Content-Type": "application/json",
    Network: useStore.getState().network,
  };

  let response: any;

  if (method === APIMethod.POST_FILE) {
    headers["Content-Type"] = "multipart/form-data";
    const form = new FormData();
    form.append("file", body.Filename);
    response = await axios
      .post(`${endpoint}`, form, { headers })
      .then((res: any) => res)
      .catch((error: any) => {
        console.log("E_REQUEST =>", error);
        if (
          error.response?.status === 401 &&
          (error.response?.data?.includes("unauthorized") ||
            error.response?.data?.includes("invalid token"))
        ) {
          console.log(error.response);
        }
        throw error;
      });
  } else {
    response = await axios({
      headers,
      method,
      url: `${endpoint}`,
      ...(method === APIMethod.GET ? { params: body } : { data: body }),
    })
      .then((res: any) => res)
      .catch((error: any) => {
        console.log("E_REQUEST =>", error);
        if (
          (error.response?.status === 401 &&
            error.response?.data?.message === "unauthorized") ||
          error.response?.data?.includes("unauthorized") ||
          error.response?.data?.includes("invalid token")
        ) {
          console.log(error.response);
        }
        throw error;
      });
  }

  return response;
};
