#!/usr/bin/env bash
# tiny test.sh — create a tmux window with N panes running nc and start server with `go run .`
NUM=${1:-3}
HOST=${2:-0.0.0.0}
PORT=${3:-6942}

# start server in the current tmux pane (if inside tmux)
if [ -n "${TMUX:-}" ]; then
  tmux send-keys -t "$TMUX_PANE" 'go run .' C-m
  sleep 0.5
fi

WIN="chat-test-$(date +%s)"

if [ -n "${TMUX:-}" ]; then
  NEW_WIN=$(tmux new-window -d -n "$WIN" -P -F '#{session_name}:#{window_index}')
else
  S="chat-test-$$"
  tmux new-session -d -s "$S" -n "$WIN"
  NEW_WIN="${S}:${WIN}"
  tmux attach-session -t "$S" &
  sleep 0.1
fi

for i in $(seq 2 "$NUM"); do tmux split-window -t "$NEW_WIN"; sleep 0.02; done
tmux select-layout -t "$NEW_WIN" tiled

for p in $(tmux list-panes -t "$NEW_WIN" -F '#{pane_index}'); do
  tmux send-keys -t "${NEW_WIN}.${p}" "nc $HOST $PORT" C-m
done

[ -n "${TMUX:-}" ] && tmux select-window -t "$NEW_WIN"
