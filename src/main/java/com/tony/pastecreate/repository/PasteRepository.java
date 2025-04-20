package com.tony.pastecreate.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import com.tony.pastecreate.model.PasteEntity;

@Repository
public interface PasteRepository extends JpaRepository<PasteEntity, Long> {
    // Custom query methods (if any)
}
