local channel, args, noWait = ...

if vim.bo.buftype == 'terminal' then vim.api.nvim_win_hide(0) end

-- nvim-style arguments:
--   +N        jump to line N
--   +/pat     search for pattern
--   +cmd      run Ex command
-- Such a "+" argument applies to the first file edited, like nvim does.
local files, cmds = {}, {}
for _, v in ipairs(args) do
  if v:sub(1, 1) == '+' and #v > 1 then
    table.insert(cmds, v:sub(2))
  else
    table.insert(files, v)
  end
end

if #files == 0 then
  for _, cmd in ipairs(cmds) do
    vim.cmd(cmd)
  end
end

local first = true
for _, v in ipairs(files) do
  vim.cmd.edit(v)
  if first then
    for _, cmd in ipairs(cmds) do
      if cmd:match('^%d+$') then
        -- jump to line N
        vim.api.nvim_win_set_cursor(0, { math.max(1, tonumber(cmd)), 0 })
      elseif cmd:sub(1, 1) == '/' then
        vim.fn.search(cmd:sub(2))
      else
        vim.cmd(cmd)
      end
    end
    first = false
  end
end

local single = #files == 1 and not noWait
if single then
  vim.api.nvim_create_autocmd({ 'BufWritePost', 'BufDelete' }, {
    buffer = 0,
    once = true,
    callback = function() vim.rpcnotify(channel, 'stop') end,
  })
else
  vim.rpcnotify(channel, 'stop')
end