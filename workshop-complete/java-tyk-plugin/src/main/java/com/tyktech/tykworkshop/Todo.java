package com.tyktech.tykworkshop;

import com.google.gson.annotations.SerializedName;
import java.sql.Timestamp;

public class Todo {
    @SerializedName("id") public String ID;

    @SerializedName("user") public String User;

    @SerializedName("todo") public String Todo;

    @SerializedName("complete") public Boolean Complete;

    @SerializedName("created_at") public java.sql.Timestamp CreatedAt;
}