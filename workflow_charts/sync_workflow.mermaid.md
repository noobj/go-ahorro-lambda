flowchart LR
    A(((Client))) -->|Sync request| SyncR(Sync Receiver)
    SQS -->|Pop handle request| D(Sync Handler)
    D -->|Update task status to 'Done'| Dynamo
    SyncR --> IF{Google OAuth Ok}
    IF -->|Success - return task id|A
    IF -->|Success - request sync task|SyncTask
    IF -->|Failed - redirect to auth page|SyncCB(Sync Callback)
    SyncCB -->|return task id|A
    SyncCB -->SyncTask
    SyncTask -->|Insert sync task status 'Pending'|Dynamo
    SyncTask -->|Push handle request| SQS
    A -->|Long polling with task id| G(Sync Notifier)
    G -->|Update with new entries when task done| A
    Dynamo -->|Check task status| G