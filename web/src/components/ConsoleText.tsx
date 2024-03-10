import React, { useEffect, useRef, useState } from "react";

const addCharacter = (text: string, i: number) => {
  let startAt = i;
  switch (text.charAt(i)) {
    // images are formatted as ![alt text](image link)
    case "!":
      if (text.charAt(i + 1) !== "[") {
        return {
          count: 1,
          text: text.charAt(i),
        };
      }
      let imageText = "";
      i++;
      while (text.charAt(i) !== "(") {
        imageText += text.charAt(i);
        i++;
      }
      i++;
      let imageLink = "";
      while (text.charAt(i) !== ")") {
        imageLink += text.charAt(i);
        i++;
      }
      i++;
      imageText = imageText.substring(1, imageText.length - 1); // strip [ and ]
      return {
        count: i - startAt,
        text: `<img src="${imageLink}" title="${imageText}" alt="${imageText}" class="w-full md:w-1/2 m-auto h-auto" />`,
      };
    case "#":
      i++;
      while (text.charAt(i) === " ") {
        i++;
      }
      let headerText = "";
      while (text.charAt(i) !== "\n") {
        headerText += text.charAt(i);
        i++;
      }
      i++;
      return {
        count: i - startAt,
        text: `<h1 class="text-xl font-bold">${headerText}</h1>`,
      };
    case "\n":
      return {
        count: 1,
        text: "<br/><br/>",
      };
    case "*":
      let boldText = "";
      i++;
      while (text.charAt(i) !== "*") {
        boldText += text.charAt(i);
        i++;
      }
      i++;
      return {
        count: i - startAt,
        text: `<b>${boldText}</b>`,
      };
    case "_":
      let underlineText = "";
      i++;
      while (text.charAt(i) !== "_") {
        underlineText += text.charAt(i);
        i++;
      }
      i++;
      return {
        count: i - startAt,
        text: `<u>${underlineText}</u>`,
      };
    case "`":
      let codeText = "";
      i++;
      while (text.charAt(i) !== "`") {
        codeText += text.charAt(i);
        i++;
      }
      i++;
      return {
        count: i - startAt,
        text: `<code>${codeText}</code>`,
      };
    case "/":
      let italicText = "";
      i++;
      while (text.charAt(i) !== "/") {
        italicText += text.charAt(i);
        i++;
      }
      i++;
      return {
        count: i - startAt,
        text: `<i>${italicText}</i>`,
      };
    default:
      return {
        count: 1,
        text: text.charAt(i),
      };
  }
};

const ConsoleText: React.FC<{ text: string }> = ({ text }) => {
  const [charByCharInterval, setCharByCharInterval] = useState<number | undefined>(undefined);
  const consoleTextRef = useRef<HTMLSpanElement>(null);
  const [rendered, setRendered] = useState("");

  useEffect(() => {
    if (text && consoleTextRef.current) {
      const consoleText = consoleTextRef.current;

      consoleText.innerHTML = "";
      let i = 0;
      // print character by character with a delay
      const intervalId = setInterval(() => {
        const added = addCharacter(text, i);

        i += added.count;
        consoleText.innerHTML += added.text;

        if (i >= text.length) {
          setRendered(text);
          clearInterval(intervalId);
        }
      }, 20);

      setCharByCharInterval(intervalId as unknown as number);

      return () => {
        clearInterval(intervalId);
      };
    }
  }, [text]);

  // can tap to skip the text animation
  const finishAnimation = () => {
    if (consoleTextRef.current && rendered !== text) {
      clearInterval(charByCharInterval);

      const actualText = text;

      const consoleText = consoleTextRef.current;
      let inc = 1;
      let html = "";
      for (let i = 0; i < text.length; i += inc) {
        const { count, text } = addCharacter(actualText, i);
        html += text;
        inc = count;
      }
      consoleText.innerHTML = html;
      setRendered(actualText);
    }
  };

  return (
    <div className="flex flex-grow overflow-y-scroll overscroll-contain p-4">
      <span ref={consoleTextRef} onClick={finishAnimation} className="w-full font-mono text-harlequin-700"></span>
    </div>
  );
};

export default ConsoleText;
