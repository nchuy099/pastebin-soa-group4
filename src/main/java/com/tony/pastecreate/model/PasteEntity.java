package com.tony.pastecreate.model;

import jakarta.persistence.*;

@Entity
@Table(name = "paste")
public class PasteEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String content;

    // Getters and setters (or use Lombok's @Data)
}
