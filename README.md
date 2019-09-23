# packer-provisioner-lumberjack
A Packer plugin that truncates logs to 0 bytes as a means of sanitization.

## Configuration Options
|option|default|description|
|---|---|---|
|`print_only`|`false`|Print (without truncating) the effected files.|
|`prevent_sudo`|`false`|Do, or do not, use `sudo`|
|`exclude_paths`|`none`|Paths to explicitly exclude from truncation.|
|`base_command`|`find / -name "*.log"`|The initial command onto which "not if path" logic is appended|

## Assumptions
  - `base_command` must accept GNU Find (of `findutils`) style logice (` -a -not -path <path> `).
  - Your source AMI must contain GNU Find or which ever `base_command` you've configured.
  - If you've configured `prevent_sudo` to `false`, `sudo` is available.

## Example usage
See [test/template.json](test/template.json):

```json
{
  "builders": [
    {
      "type": "docker",
      "image": "fedora/tools",
      "discard": true
    }
  ],
  "provisioners": [
    {
      "type": "lumberjack",
      "print_only": false,
      "exclude_paths": [
        "/tmp"
      ],
      "prevent_sudo": true
    }
  ]
}

```

## `make`
The default make target runs `go build`, copies the plugin binary to `~/.packer.d/plugins`, executes `packer build test/template.json` (with `PACKER_LOG=1` debug logging on), and `cat`s out packer.log.

### Why 'lumberjack' (asked no-one)?
It "truncates" logs. It TRUNKactes LOGS... ugh.

### Reference
A number of `packer` plugins and `godocs` were referenced to write this one:
  - github.com/rgl/packer-provisioner-windows-update
  - github.com/vultr/packer-builder-vultr
  - github.com/SwampDragons/packer-provisioner-comment
