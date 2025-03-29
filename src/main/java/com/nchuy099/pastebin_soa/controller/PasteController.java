package com.nchuy099.pastebin_soa.controller;


import com.nchuy099.pastebin_soa.repository.projection.MonthlyStatsProjection;
import com.nchuy099.pastebin_soa.dto.response.ResponseData;
import com.nchuy099.pastebin_soa.service.PasteService;
import lombok.RequiredArgsConstructor;

import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import org.springframework.web.bind.annotation.GetMapping;

import org.springframework.web.bind.annotation.*;

import java.time.YearMonth;

@RestController
@RequestMapping("/api/paste")
@RequiredArgsConstructor
public class PasteController {

    private final PasteService pasteService;

    @GetMapping("/stats")
    public ResponseEntity<ResponseData<MonthlyStatsProjection>> getMonthlyStats(
            @RequestParam(value = "month") @DateTimeFormat(pattern = "yyyy-MM") YearMonth yearMonth) throws Exception {
        return ResponseEntity.ok(new ResponseData<>(HttpStatus.OK.value(), "Get monthly stats successfully", pasteService.getMonthlyStats(yearMonth)));

    }
}