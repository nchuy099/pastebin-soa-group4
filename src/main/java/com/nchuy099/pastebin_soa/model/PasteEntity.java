package com.nchuy099.pastebin_soa.model;

import com.nchuy099.pastebin_soa.common.Status;
import com.nchuy099.pastebin_soa.common.Visibility;
import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.hibernate.annotations.CreationTimestamp;

import java.util.Date;

@Setter
@Getter
@Entity
@Table(name = "paste")
public class PasteEntity {
    @Id
    @Column(name = "id", length = 10)
    private String id;

    @Column(name = "content", columnDefinition = "TEXT", nullable = false)
    private String content;

    @Column(name = "title", length = 255)
    private String title = "Untitled";

    @Column(name = "language", length = 50)
    private String language = "text";

    @Column(name = "created_at")
    @Temporal(TemporalType.TIMESTAMP)
    @CreationTimestamp
    private Date createdAt;

    @Column(name = "expires_at")
    @Temporal(TemporalType.TIMESTAMP)
    private Date expiresAt;

    @Column(name = "views")
    private int views = 0;

    @Enumerated(EnumType.STRING)
    @Column(name = "visibility")
    private Visibility visibility = Visibility.PUBLIC;

    @Enumerated(EnumType.STRING)
    @Column(name = "status")
    private Status status = Status.ACTIVE;
}
