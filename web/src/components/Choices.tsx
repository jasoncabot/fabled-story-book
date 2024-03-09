import React from "react";

export interface Choice {
  text: string;
  code: string;
}

const Choices: React.FC<{ choices: Choice[]; onChoiceSelected: (c: Choice) => void }> = ({ choices, onChoiceSelected }) => {
  const buttons = choices.map((choice) => {
    return (
      <button
        type="button"
        className="min-h-12 min-w-48 rounded-lg border border-harlequin-700 bg-slate-900 px-4 py-2 font-mono text-xs font-medium text-harlequin-700 hover:bg-harlequin-900 hover:text-harlequin-400"
        onClick={(e) => {
          e.preventDefault();
          onChoiceSelected(choice);
        }}
      >
        {choice.text}
      </button>
    );
  });

  return (
    <div className="flex inline-flex shrink-0 flex-row space-x-4 overflow-x-auto rounded-md p-4" role="group">
      {buttons}
    </div>
  );
};

export default Choices;
