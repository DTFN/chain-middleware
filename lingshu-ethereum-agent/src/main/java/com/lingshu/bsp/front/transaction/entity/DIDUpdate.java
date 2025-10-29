/*
 * Copyright 2014-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.lingshu.bsp.front.transaction.entity;

import lombok.Data;

/**
 * transHandleWithSign interface parameter.
 * handle transactions of deploy/call contract
 * v1.3.0+ default with sign
 */
@Data
public class DIDUpdate {
    private Integer ledgerId;
    private String address;
    private String did;
    private String didDoc;
}
