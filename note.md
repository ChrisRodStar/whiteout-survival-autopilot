# ToDo: Bot Farm Architecture Implementation

- [ ] bravo - does not parse user

## 🔁 Stage 4: Task Execution and TTL
- [ ] Process tasks sequentially, one profile at a time
- [ ] Assign TTL to each task
- [ ] Mark completed tasks by TTL in profile context
- [ ] Skip/reschedule tasks with expired TTL

## ⚖️ Stage 5: Profile Switching and Priorities
- [ ] Check for priority tasks from other profiles on the device
- [ ] Switch to profile with higher priority task when necessary (after completing current task)

## 🧠 Stage 6: Global Analyzer
- [ ] Launch separate goroutine for global state analysis
- [ ] Periodically check statuses of all profiles
- [ ] Generate urgent tasks when critical events are detected
- [ ] Add priority tasks to the queue of corresponding `gamer.id`

## 🧩 Additional (optional)
- [ ] Add logging for actions and switches
- [ ] Display profile and queue status (e.g., via debug endpoint)
- [ ] Limit queue length per profile
- [ ] Task scheduler by priority, not just FIFO
