package com.lingshu.bsp.front.util;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.ToString;
import lombok.extern.log4j.Log4j2;
import lombok.extern.slf4j.Slf4j;

import java.io.*;
import java.util.concurrent.TimeUnit;

/**
 * @author: zehao.song
 */
@Slf4j
@ToString
public class JavaCmdExecutor {
    public static final long DEFAULT_EXEC_TIMEOUT = 60_60_1000L;

    public static ExecResult executeCommand(String command, long timeout) {
        Process process = null;
        InputStream pIn = null;
        InputStream pErr = null;
        StreamGobbler outputGobbler = null;
        StreamGobbler errorGobbler = null;
        try {
            log.info("exec command:[{}]", command);
            String[] commandArray = { "/bin/bash", "-c", command };
            process = Runtime.getRuntime().exec(commandArray);
            final Process p = process;

            // close process's output stream.
            closeQuietly(p.getOutputStream());

            pIn = process.getInputStream();
            outputGobbler = new StreamGobbler(pIn, "OUTPUT");
            outputGobbler.start();

            pErr = process.getErrorStream();
            errorGobbler = new StreamGobbler(pErr, "ERROR");
            errorGobbler.start();
            long newTimeout = timeout <= 0 ? DEFAULT_EXEC_TIMEOUT : timeout;

            p.waitFor(newTimeout, TimeUnit.MILLISECONDS);
            int exitCode = p.exitValue();

            if (exitCode == 0) {
                log.info("Exec command success: code:[{}], OUTPUT:\n[{}]",
                        exitCode, outputGobbler.getContent());
            } else {
                log.warn("Exec command code not zero: code:[{}], OUTPUT:\n[{}],\nERROR:\n[{}]",
                        exitCode, outputGobbler.getContent(), errorGobbler.getContent());
            }

            return new ExecResult(exitCode, outputGobbler.getContent());
        } catch (IOException ex) {
            String errorMessage = "The command [" + command + "] execute failed.";
            log.error(errorMessage, errorMessage);
            return new ExecResult(-1, null);
        } catch (InterruptedException ex) {
            String errorMessage = "The command [" + command + "] did not complete due to an interrupted error.";
            log.error(errorMessage, ex);
            return new ExecResult(-1, errorMessage);
        } finally {
            if (pIn != null) {
                closeQuietly(pIn);
                if (outputGobbler != null && !outputGobbler.isInterrupted()) {
                    outputGobbler.interrupt();
                }
            }
            if (pErr != null) {
                closeQuietly(pErr);
                if (errorGobbler != null && !errorGobbler.isInterrupted()) {
                    errorGobbler.interrupt();
                }
            }
            if (process != null) {
                process.destroy();
            }
        }
    }

    private static void closeQuietly(Closeable c) {
        try {
            if (c != null) {
                c.close();
            }
        } catch (IOException e) {
            log.error("Exception occurred when closeQuietly!!! ", e);
        }
    }

    @Data
    @ToString
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ExecResult {
        private int exitCode;
        private String executeOut;

        public static ExecResult build (int exitCode,String msg){
            return new ExecResult(exitCode,msg);
        }

        public boolean success(){
            return exitCode == 0;
        }
        public boolean failed(){
            return !success();
        }
    }

    @Slf4j
    public static class StreamGobbler extends Thread {
        private InputStream inputStream;
        private String streamType;
        private StringBuilder buf;
        private boolean isStopped = false;

        /**
         * @param inputStream the InputStream to be consumed
         * @param streamType  the stream type (should be OUTPUT or ERROR)
         */
        public StreamGobbler(final InputStream inputStream, final String streamType) {
            this.inputStream = inputStream;
            this.streamType = streamType;
            this.buf = new StringBuilder();
            this.isStopped = false;
        }

        /**
         * Consumes the output from the input stream and displays the lines consumed
         * if configured to do so.
         */
        @Override
        public void run() {
            try {
                // 默认编码为UTF-8，这里设置编码为GBK，因为WIN7的编码为GBK
                InputStreamReader inputStreamReader = new InputStreamReader(inputStream, "UTF-8");
                BufferedReader bufferedReader = new BufferedReader(inputStreamReader);
                String line = null;
                while ((line = bufferedReader.readLine()) != null) {
                    this.buf.append(line + "\n");
                }
            } catch (IOException ex) {
                log.error("Failed to successfully consume and display the input stream of type {}.", streamType, ex);
            } finally {
                this.isStopped = true;
                synchronized (this) {
                    notifyAll();
                }
            }
        }

        public String getContent() {
            if (!this.isStopped) {
                synchronized (this) {
                    try {
                        wait();
                    } catch (InterruptedException ignore) {
                        ignore.printStackTrace();
                        Thread.currentThread().interrupt();
                    }
                }
            }
            return this.buf.toString();
        }
    }
}
