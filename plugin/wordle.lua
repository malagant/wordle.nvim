-- wordle.nvim — Play Wordle in Neovim
-- Usage: :Wordle [de|en] [random]

local M = {}

-- Find the wordle-nvim binary
local function find_binary()
  -- Check if installed via go install
  local go_bin = vim.fn.expand("$HOME/go/bin/wordle-nvim")
  if vim.fn.executable(go_bin) == 1 then
    return go_bin
  end

  -- Check plugin directory (for development / goreleaser)
  local plugin_dir = vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":h:h")
  local local_bin = plugin_dir .. "/wordle-nvim"
  if vim.fn.executable(local_bin) == 1 then
    return local_bin
  end

  -- Check PATH
  if vim.fn.executable("wordle-nvim") == 1 then
    return "wordle-nvim"
  end

  return nil
end

function M.play(args)
  local binary = find_binary()
  if not binary then
    vim.notify("wordle-nvim binary not found! Run: go install github.com/malagant/wordle-nvim/cmd/wordle@latest", vim.log.levels.ERROR)
    return
  end

  local cmd = binary
  for _, arg in ipairs(args) do
    cmd = cmd .. " " .. arg
  end

  -- Open in a new terminal buffer
  vim.cmd("botright new")
  vim.cmd("resize 25")
  vim.fn.termopen(cmd, {
    on_exit = function(_, code, _)
      if code == 0 then
        vim.schedule(function()
          -- Close the buffer when the game ends cleanly
          local buf = vim.api.nvim_get_current_buf()
          if vim.api.nvim_buf_is_valid(buf) then
            vim.api.nvim_buf_delete(buf, { force = true })
          end
        end)
      end
    end,
  })
  vim.cmd("startinsert")
end

-- Register the :Wordle command
vim.api.nvim_create_user_command("Wordle", function(opts)
  M.play(opts.fargs)
end, {
  nargs = "*",
  desc = "Play Wordle in the terminal",
  complete = function()
    return { "en", "de", "random" }
  end,
})

return M
