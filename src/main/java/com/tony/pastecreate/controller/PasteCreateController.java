package com.tony.pastecreate.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import com.tony.pastecreate.model.PasteEntity;
import com.tony.pastecreate.service.PasteCreateService;

@RestController
@RequestMapping("/api/paste")
public class PasteCreateController {

    private final PasteCreateService pasteCreateService;

    @Autowired
    public PasteCreateController(PasteCreateService pasteCreateService) {
        this.pasteCreateService = pasteCreateService;
    }

    @PostMapping
    public ResponseEntity<PasteEntity> createPaste(@RequestBody PasteEntity pasteEntity) {
        PasteEntity createdPaste = pasteCreateService.createPaste(pasteEntity);
        return ResponseEntity.ok(createdPaste);
    }
}