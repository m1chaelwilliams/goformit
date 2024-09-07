# GoFormIt

A runtime for terminal UI forms

>[!warning] This product is in alpha. New features will be breaking until the project reached 1.0.

## Features

- Interactive Terminal User Interface (TUI)
- JSON support
- Non-Linear Forms
- State management
- 🚧 Experimental go package

## Installation & Usage

As of right now, the project requires building from source. In order to do so, you will need:

- Golang

First, clone the repository:

```shell
git clone https://github.com/m1chaelwilliams/goformit
```

Next, `cd` into the directory and build the project:

```shell
cd goformit
go build
```

Ensure everything works by trying out one of the test forms:

```shell
./goformit -i testforms/sample_form.json -o out.json (-v)
```

- `i` is the input filepath
- `o` is the output filepath
- `v` is the optional verbose flag, which dumps logging to stdout after the program finishes.

## Making a Form

The form builder TUI is currently under construction. Until it is useable, creating a form
will require creating the JSON file from scratch.

A form has **3** root fields:

- prompts
- first_prompt
- vars

### "prompts"

A form is made up of smaller pieces called "prompts". Here is what a sample prompt
would look like in JSON:

```json
"my_prompt_id": {
  "id": "my_prompt_id",
  "group": "my_group",
  "type": "selection",
  "title": "This is a prompt",
  "choices": [
    "Option A",
    "Option B"
  ],
  "next": {
    "Option A": "next_prompt_id",
    "_": "[[end]]"
  }
},
```

There is quite a bit of information here. Let's break it down:

- id (required): a unique identifier for the prompt
- group (optional): a unique group id that the prompt response will end up in in the ouput
- type (required): the type of the prompt (input, selection, checkbox)
- title (required): the display title of the prompt
- choices (required for selection and multiselection): the display options to choose from
- next (required): a mapping of response -> next prompt ([[end]] signifies the end of the program)

### "first_prompt"

This is the unique id of the first prompt (e.g. entry point) of the form. This is a required field.

### "vars" (Experimental)

Perhaps the most complex part of the runtime, the "vars" are a way to manage the form's state without having access to the codebase.

#### Example Usage

```json
"vars": {
    "group_id": "default"
}
"prompts": {
    "my_prompt_id": {
      "id": "my_prompt_id",
      "group": "[[group_id]]",
      "bind_submit": "[[group_id]]"
      "type": "input",
      "title": "Choose which group I end up in",
      "next": {
        "_": "next_prompt"
      }
    }
}
```

In this example, a variable called `group_id` is declared with a default value of "default". Then, in the prompt, `group_id` is used with the special syntax `[[var]]` inside of "group" and "bind_submit". What this means is that 1. The group of `my_prompt_id` will be the value of `group_id`; and 2. on submit, `group_id` will be set to the prompt's response.

This feature is highly experimental and in its early stages. Any feedback/suggestions are appreciated.

## Credit

The TUI toolkit used for this project is [Bubbletea](https://github.com/charmbracelet/bubbletea).

## License

This project is licensed under [MIT](./LICENSE).

## Contributing

Any contributions are welcome! Please take a look at the [roadmap](./ROADMAP.md) to see what is
most needed.

## Contributors

Michael Williams - Creator
[Support me on Ko-fi](https://ko-fi.com/codingwithsphere)
