import { Cog6ToothIcon } from "@heroicons/react/24/outline";
import React, { useContext, useEffect, useState } from "react";
import Choices, { Choice } from "./components/Choices";
import ConsoleText from "./components/ConsoleText";
import { wasmInterop } from "./jabl";
import jablEngine from "./jabl.wasm";
import { sources } from "./sources";

const defaultEntrypoint = "entrypoint.jabl";

export interface Result {
  output: string;
  choices: Choice[];
  transition: string;
}

const App: React.FC = () => {
  return (
    <JablProvider>
      <GameProvider>
        <Content />
      </GameProvider>
    </JablProvider>
  );
};

const Content: React.FC = () => {
  const [settingsHidden, setSettingsHidden] = useState(true);
  const { transitionCount, section, source, resetProgress, changeStory, transition } = useContext(GameContext);
  const { execSection, evalCode } = useContext(WasmContext);

  const [result, setResult] = useState<Result | undefined>(undefined);
  const [error, setError] = useState<string | undefined>(undefined);

  const update = (result: string | undefined, error: string | undefined) => {
    if (error) {
      setResult(undefined);
      setError(error);
    } else if (result) {
      const parsedResult = JSON.parse(result) as Result;
      if (parsedResult.transition.length > 0) {
        transition(parsedResult.transition);
      } else {
        setResult(parsedResult);
        setError(undefined);
      }
    }
  };

  useEffect(() => {
    if (section && source) {
      execSection(section, update);
    } else if (source) {
      const entrypoint = sources.find((s) => s.id === source)?.entrypoint || defaultEntrypoint;
      execSection(entrypoint, update);
    } else {
      const choices = sources.map((source) => `choice("${source.name}", { set("system:source", ${source.id}) goto("${source.entrypoint}")})`);
      evalCode(
        `{
        print("Welcome to the game!")
        print("Which book would you like to play?")
        ${choices.join("\n")}
      }`,
        update,
      );
    }
  }, [transitionCount, section, source, evalCode, execSection]);

  return (
    <div className="flex h-[100svh] flex-col bg-slate-900">
      <div className="flex flex-shrink flex-row justify-end">
        <ul className={"options absolute mt-8" + (settingsHidden ? " hidden" : "")}>
          <li>
            <a
              href="#"
              onClick={() => {
                setSettingsHidden(true);
                resetProgress();
              }}
              className="block bg-slate-800 px-4 py-2 font-mono text-harlequin-700 hover:bg-slate-900 hover:text-harlequin-400"
            >
              Reset progress
            </a>
          </li>
          <li>
            <a
              href="#"
              onClick={() => {
                setSettingsHidden(true);
                changeStory();
              }}
              className="block bg-slate-800 px-4 py-2 font-mono text-harlequin-700 hover:bg-slate-900 hover:text-harlequin-400"
            >
              Change story
            </a>
          </li>
        </ul>

        <button onClick={() => setSettingsHidden(!settingsHidden)}>
          <Cog6ToothIcon className="m-2 h-8 w-8 text-harlequin-700" />
        </button>
      </div>

      {error && (
        <div className="mx-4 mb-4 border-b border-t border-red-500 bg-red-100 px-4 py-3 text-red-700" role="alert">
          {error}
        </div>
      )}

      <ConsoleText text={result?.output || ""}></ConsoleText>

      <Choices choices={result?.choices ?? []} onChoiceSelected={(choice) => evalCode(choice.code, update)} />
    </div>
  );
};

interface GameContextProps {
  transitionCount: number;
  section: string | undefined;
  source: string | undefined;
  resetProgress: () => void;
  changeStory: () => void;
  transition: (transition: string) => void;
}

const GameContext = React.createContext<GameContextProps>({
  transitionCount: 0,
  section: undefined,
  source: undefined,
  resetProgress: () => {},
  changeStory: () => {},
  transition: (_) => {},
});

const GameProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [transitionCount, setTransitionCount] = useState(0);
  const [source, setSource] = useState(wasmInterop.bookStorage.getItem("system:source"));
  const [section, setSection] = useState(wasmInterop.bookStorage.getItem("section"));

  const resetProgress = () => {
    const entryPoint = sources.find((s) => s.id === source)?.entrypoint || defaultEntrypoint;
    wasmInterop.bookStorage.setItem("section", entryPoint);
    wasmInterop.bookStorage.clearCurrent();
    setSection(entryPoint);
    setSource(wasmInterop.bookStorage.getItem("system:source"));
    setTransitionCount(transitionCount + 1);
  };

  const changeStory = () => {
    wasmInterop.bookStorage.setItem("system:source", undefined);
    setSection(undefined);
    setSource(undefined);
    setTransitionCount(0);
  };

  const transition = (transition: string) => {
    wasmInterop.bookStorage.setItem("section", transition);
    setSection(transition);
    setSource(wasmInterop.bookStorage.getItem("system:source"));
    setTransitionCount(transitionCount + 1);
  };

  return <GameContext.Provider value={{ transitionCount, section, source, resetProgress, changeStory, transition }}>{children}</GameContext.Provider>;
};

interface WasmContextProps {
  execSection: (section: string, callback: (code: string, error: string) => void) => void;
  evalCode: (code: string, callback: (result: string, error: string) => void) => void;
}

const WasmContext = React.createContext<WasmContextProps>({
  execSection: () => {},
  evalCode: () => {},
});

const JablProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [wasm, setWasm] = useState<WebAssembly.Instance | null>(null);

  useEffect(() => {
    if (!WebAssembly.instantiateStreaming) {
      WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
      };
    }

    const go = new Go();
    WebAssembly.instantiateStreaming(fetch(jablEngine), go.importObject).then(({ instance }) => {
      go.run(instance);
      setWasm(instance);
    });
    Object.assign(window, wasmInterop);
  }, []);

  const execSection = (section: string, callback: (code: string, error: string) => void) => {
    if (wasm) {
      (window as any).execSection(section, callback);
    }
  };

  const evalCode = (code: string, callback: (result: string, error: string) => void) => {
    if (wasm) {
      (window as any).evalCode(code, callback);
    }
  };

  return <WasmContext.Provider value={{ execSection, evalCode }}>{children}</WasmContext.Provider>;
};

export default App;
