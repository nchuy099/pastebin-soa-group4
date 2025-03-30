package com.tony.pastecreate.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import com.tony.pastecreate.model.PasteEntity;

public interface PasteRepository extends JpaRepository<PasteEntity, Long> {
}
