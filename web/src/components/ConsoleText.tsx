import React, { useEffect, useRef, useState } from "react";

const ConsoleText: React.FC<{ text: string }> = ({ text }) => {
  const consoleTextRef = useRef<HTMLSpanElement>(null);
  const [intervalId, setIntervalId] = useState<number | undefined>(undefined);

  useEffect(() => {
    if (consoleTextRef.current) {
      const consoleText = consoleTextRef.current;
      consoleText.innerHTML = "";
      let textIndex = 0;
      const intervalId = setInterval(() => {
        if (text && textIndex < text.length) {
          switch (text.charAt(textIndex)) {
            case "\n":
              consoleText.innerHTML += "<br/><br/>";
              break;
            case "*":
              let boldText = "";
              textIndex++;
              while (text.charAt(textIndex) !== "*") {
                boldText += text.charAt(textIndex);
                textIndex++;
              }
              consoleText.innerHTML += `<b>${boldText}</b>`;
              break;
            case "_":
              let underlineText = "";
              textIndex++;
              while (text.charAt(textIndex) !== "_") {
                underlineText += text.charAt(textIndex);
                textIndex++;
              }
              consoleText.innerHTML += `<u>${underlineText}</u>`;
              break;
            case "`":
              let codeText = "";
              textIndex++;
              while (text.charAt(textIndex) !== "`") {
                codeText += text.charAt(textIndex);
                textIndex++;
              }
              consoleText.innerHTML += `<code>${codeText}</code>`;
              break;
            case "/":
              let italicText = "";
              textIndex++;
              while (text.charAt(textIndex) !== "/") {
                italicText += text.charAt(textIndex);
                textIndex++;
              }
              consoleText.innerHTML += `<i>${italicText}</i>`;
              break;
            default:
              consoleText.innerHTML += text.charAt(textIndex);
              break;
          }

          // scroll consoleText to the bottom of what it is displaying
          consoleText.parentElement!.scrollTop = consoleText.parentElement!.scrollHeight;
          textIndex++;
        } else {
          clearInterval(intervalId);
        }
      }, 10);
      setIntervalId(intervalId as unknown as number);
      return () => {
        clearInterval(intervalId);
      };
    }
  }, [text]);

  // can tap to skip the text animation
  const finishAnimation = () => {
    if (consoleTextRef.current) {
      const consoleText = consoleTextRef.current;
      clearInterval(intervalId);
      consoleText.innerHTML = "";
      let html = text.replace(/\/([^/]+?)\//g, "<i>$1</i>");

      html = html.replace(/\n/g, "<br/><br/>");
      html = html.replace(/\*([^\*]*?)\*/g, "<b>$1</b>");
      html = html.replace(/_([^_]+?)_/g, "<u>$1</u>");
      html = html.replace(/`([^`]+)`/g, "<code>$1</code>");

      consoleText.innerHTML = html;
      consoleText.parentElement!.scrollTop = consoleText.parentElement!.scrollHeight;
    }
  };

  return (
    <div className="flex flex-grow overflow-y-scroll overscroll-contain p-4">
      <span ref={consoleTextRef} onClick={finishAnimation} className="font-mono text-harlequin-700"></span>;
    </div>
  );
};

export default ConsoleText;
