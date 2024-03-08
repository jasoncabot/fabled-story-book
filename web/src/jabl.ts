import { authFunctions, availableSources } from "./sources";

export type WasmFunctions = {
  bookStorage: {
    getItem: (key: string) => string | undefined;
    setItem: (key: string, value: string | undefined) => void;
    clearCurrent: () => void;
  };
  loadSection: (identifier: string, callback: (code: string | undefined, error: string | undefined) => void) => void;
};

export const wasmInterop: WasmFunctions = {
  bookStorage: {
    getItem: (key: string) => {
      const sourceId = localStorage.getItem("system:source");
      const prefix = sourceId && key != "system:source" ? `fsb:${sourceId}:` : "";
      return localStorage.getItem(prefix + key) ?? undefined;
    },
    setItem: (key: string, value: string | undefined) => {
      const sourceId = localStorage.getItem("system:source");
      const prefix = sourceId && key != "system:source" ? `fsb:${sourceId}:` : "";
      if (!value) {
        localStorage.removeItem(prefix + key);
      } else {
        localStorage.setItem(prefix + key, value);
      }
    },
    clearCurrent: () => {
      const sourceId = localStorage.getItem("system:source");
      if (sourceId) {
        const prefix = `fsb:${sourceId}:`;
        for (let i = localStorage.length - 1; i >= 0; i--) {
          const key = localStorage.key(i);
          if (key && key.startsWith(prefix)) {
            localStorage.removeItem(key);
          }
        }
      }
    },
  },
  loadSection: async (identifier: string, callback: (code: string | undefined, error: string | undefined) => void) => {
    const sourceId = localStorage.getItem("system:source") ?? "";
    const sourceURL = availableSources[sourceId];
    if (!sourceURL) {
      throw new Error("Invalid source id");
    }
    // create custom headers to add to the fetch request
    const authFn = authFunctions[sourceId];
    try {
      const token = await authFn();
      const headers = new Headers();
      if (token) {
        headers.append("Authorization", `Bearer ${token}`);
      }
      const response = await fetch(`${sourceURL}${identifier}`, {
        method: "GET",
        headers: headers,
      });

      if (response.status == 401 || response.status == 403) {
        throw new Error("Unauthorised");
      } else if (response.status == 404) {
        throw new Error("Not found");
      } else if (!response.ok) {
        throw new Error("Not ok");
      } else if ((response.headers.get("content-type") ?? "").includes("text/html")) {
        throw new Error("Not a valid jabl file");
      }

      const text = await response.text();

      callback(text, undefined);
    } catch (err: Error | any) {
      callback(undefined, `Failed to fetch section ${identifier}: ${err}`);
    }
  },
};
