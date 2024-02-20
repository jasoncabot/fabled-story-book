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
    button.addEventListener("click", onChoiceSelected.bind(null, choice.code));
    choiceButtons.appendChild(button);
  });
};

const onChoiceSelected = (code) => {
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
  if (result && result.transition && result.transition.length > 0) {
    const next = exec(result.transition)
      .then((a) => {
        render(a, null);
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

  window.loadSection = (identifier, callback) => {
    console.log("loading section " + identifier);
    fetch(
      "https://raw.githubusercontent.com/jasoncabot/fabled-story-book/main/assets/example01/" +
        identifier
    )
      .then((response) => response.text())
      .then((text) => {
        callback(text, null);
      })
      .catch((err) => {
        callback(null, err);
      });
  };

  exec("entrypoint.jabl")
    .then((result) => {
      render(result, null);
    })
    .catch((e) => {
      render(null, e);
    });
};

const exec = (sectionId) => {
  return new Promise((resolve, reject) => {
    console.log("executing section " + sectionId);
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
