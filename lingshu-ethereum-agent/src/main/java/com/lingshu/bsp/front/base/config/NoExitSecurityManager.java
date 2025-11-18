package com.lingshu.bsp.front.base.config;

import java.security.Permission;

public class NoExitSecurityManager extends SecurityManager {
        @Override
        public void checkExit(int status) {
            throw new SecurityException();
        }

        @Override
      public void checkPermission( Permission permission ) {
        if( "exitVM".equals( permission.getName() ) ) {
            throw new SecurityException();
        }
    }

}
