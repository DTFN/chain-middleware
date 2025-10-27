package com.lingshu.fabric.agent;

import lombok.Data;
import lombok.Setter;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;

import java.util.Set;


public class FabricUser implements User {

    private String name;
    private Set<String> roles;
    private String account;
    private String affiliation;
    private String mspId;
    Enrollment enrollment = null; //need access in test env.

    public void setName(String name) {
        this.name = name;
    }

    public void setMspId(String mspId) {
        this.mspId = mspId;
    }

    public void setEnrollment(Enrollment enrollment) {
        this.enrollment = enrollment;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public Set<String> getRoles() {
        return roles;
    }

    @Override
    public String getAccount() {
        return account;
    }

    @Override
    public String getAffiliation() {
        return affiliation;
    }

    @Override
    public Enrollment getEnrollment() {
        return enrollment;
    }

    @Override
    public String getMspId() {
        return mspId;
    }
}
