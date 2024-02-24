let textIndex = 0;
let intervalId = null;

const renderText = (text) => {
  const consoleText = document.querySelector(".console-text");
  if (intervalId) {
    consoleText.innerHTML = "";
    clearInterval(intervalId);
    textIndex = 0;
  }

  intervalId = setInterval(() => {
    if (text && textIndex < text.length) {
      if (text.charAt(textIndex) === "\n") {
        consoleText.innerHTML += "<br/><br/>";
      } else {
        consoleText.innerHTML += text.charAt(textIndex);
      }
      // scroll consoleText to the bottom of what it is displaying
      consoleText.parentElement.scrollTop =
        consoleText.parentElement.scrollHeight;
      textIndex++;
    } else {
      clearInterval(intervalId);
    }
  }, 10);
};

const renderChoices = (choices) => {
  const choiceButtons = document.querySelector(".buttons");
  choiceButtons.innerHTML = "";
  (choices || []).forEach((choice) => {
    const button = document.createElement("button");
    button.classList.add("button");

    button.innerHTML = choice.text;
    button.addEventListener("click", runJABL.bind(null, choice.code));
    choiceButtons.appendChild(button);
  });
};

const runJABL = (code) => {
  jablEval(code)
    .then((result) => {
      render(result, null);
    })
    .catch((e) => {
      render(null, e);
    });
};

const render = (result, err) => {
  if (err) {
    renderText(err);
    return;
  }

  renderText(result.output);
  renderChoices(result.choices);
  const transition = result?.transition || "";
  if (transition.length > 0) {
    exec(transition)
      .then((response) => {
        render(response, null);
      })
      .catch((e) => {
        render(null, e);
      });
  }
};

const run = async () => {
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const go = new Go();
  const { instance } = await WebAssembly.instantiateStreaming(
    fetch("test.wasm"),
    go.importObject
  );
  go.run(instance);

  registerGlobals();

  // Each source has it's own storage so even if they use the same variables they don't collide
  if (localStorage.getItem("system:source")) {
    startSelection();
  } else {
    showSelectionChoices();
  }
};

const registerGlobals = () => {
  window.bookStorage = {
    getItem: (key) => {
      const sourceId = localStorage.getItem("system:source");
      const prefix = sourceId ? `fsb:${sourceId}:` : "";
      return localStorage.getItem(prefix + key);
    },
    setItem: (key, value) => {
      const sourceId = localStorage.getItem("system:source");
      const prefix = sourceId ? `fsb:${sourceId}:` : "";
      localStorage.setItem(prefix + key, value);
    },
  };
  window.resetProgress = () => {
    // hide the settings menu
    document.querySelector(".popover").classList.add("hidden");

    // Remove appropriate keys from local storage
    const sourceId = localStorage.getItem("system:source");
    const prefix = sourceId ? `fsb:${sourceId}:` : "";
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key.startsWith(prefix)) {
        localStorage.removeItem(key);
      }
    }
    localStorage.setItem("system:section", "entrypoint.jabl");

    // And restart the game
    startSelection();
  };
  window.changeStory = () => {
    // hide the settings menu
    document.querySelector(".popover").classList.add("hidden");

    localStorage.removeItem("system:source");
    localStorage.removeItem("system:section");

    showSelectionChoices();
  };
  window.loadSection = (identifier, callback) => {
    const sourceId = localStorage.getItem("system:source");
    const sourceURL = {
      1: "https://raw.githubusercontent.com/jasoncabot/fabled-story-book/main/assets/example01/",
    }[sourceId];
    if (!sourceURL) {
      throw new Error("Invalid source id");
    }
    fetch(sourceURL + identifier)
      .then((response) => response.text())
      .then((text) => {
        callback(text, null);
      })
      .catch((err) => {
        callback(null, err);
      });
  };
};

const startSelection = () => {
  let currentSection = localStorage.getItem("system:section");
  exec(currentSection)
    .then((result) => {
      render(result, null);
    })
    .catch((e) => {
      render(null, e);
    });
};

const showSelectionChoices = () => {
  runJABL(`{
    print("Welcome to the game!")
    print("Which book would you like to play?")
    choice("Example 1", {
      set("system:source", 1)
      goto("entrypoint.jabl")
    })
  }`);
};

const exec = (sectionId) => {
  return new Promise((resolve, reject) => {
    localStorage.setItem("system:section", sectionId);
    window.execSection(sectionId, (value, err) => {
      if (err) {
        reject(err);
      } else {
        try {
          result = JSON.parse(value);
          resolve(result);
        } catch (e) {
          reject(e);
        }
      }
    });
  });
};

const jablEval = (code) => {
  return new Promise((resolve, reject) =>
    window.evalCode(code, (value, err) => {
      if (err) {
        reject(err);
      } else {
        resolve(JSON.parse(value));
      }
    })
  );
};

const toggleSettings = () => {
  const popover = document.querySelector(".popover");
  popover.classList.toggle("hidden");
};
