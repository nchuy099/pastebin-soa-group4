package com.tony.pastecreate.controller;

import com.tony.pastecreate.model.PasteEntity;
import com.tony.pastecreate.service.PasteCreateService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/paste")
@RequiredArgsConstructor
public class PasteCreateController {
    private final PasteCreateService pasteCreateService;

    @PostMapping
    public ResponseEntity<PasteEntity> createPaste(@RequestBody PasteEntity paste) {
        PasteEntity savedPaste = pasteCreateService.createPaste(paste);
        return ResponseEntity.ok(savedPaste);
    }
}