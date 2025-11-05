local channel, files = ...

if vim.bo.buftype == 'terminal' then vim.api.nvim_win_hide(0) end

for _, v in ipairs(files) do
  vim.cmd.edit(v)
end

if #files == 1 then
  vim.api.nvim_create_autocmd({ 'BufWritePost', 'BufDelete' }, {
    buffer = 0,
    once = true,
    callback = function() vim.rpcnotify(channel, 'stop') end,
  })
else
  vim.rpcnotify(channel, 'stop')
end
