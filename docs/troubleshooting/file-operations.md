# File Operations

## Upload Failures

**Symptom:**
```
Error: failed to upload file
Error: invalid file path
Error: file too large
```

**Solutions:**

1. **File not found:**
   ```bash
   # Verify file exists
   ls -lh /path/to/file.tar.gz
   
   # Use absolute paths
   go-dci upload-file --job-id <job-id> --file /absolute/path/to/file.tar.gz
   ```

2. **Permission denied:**
   ```bash
   # Check file permissions
   chmod 644 /path/to/file.tar.gz
   ```

3. **File too large:**
   - DCI has file size limits (typically 1GB)
   - Compress large files before uploading:
     ```bash
     tar czf results.tar.gz results/
     go-dci upload-file --job-id <job-id> --file results.tar.gz --mime-type application/gzip
     ```

4. **Wrong MIME type:**
   Specify the correct MIME type:
   ```bash
   go-dci upload-file --job-id <job-id> --file results.xml --mime-type application/xml
   go-dci upload-file --job-id <job-id> --file results.json --mime-type application/json
   go-dci upload-file --job-id <job-id> --file results.tar.gz --mime-type application/gzip
   ```

## Download Issues

**Symptom:**
```
Error: failed to download file
Error: file content empty
```

**Solutions:**

1. Verify the file ID is correct:
   ```bash
   # List files for a job
   go-dci job-files --job-id <job-id>
   ```

2. Check if you have permission to access the file (must be in the same team)

3. For library users, handle the response properly:
   ```go
   data, contentType, err := client.GetFile(fileID)
   if err != nil {
       log.Fatalf("Download failed: %v", err)
   }
   
   // Save to file
   err = os.WriteFile("output.tar.gz", data, 0644)
   if err != nil {
       log.Fatalf("Failed to save file: %v", err)
   }
   ```
